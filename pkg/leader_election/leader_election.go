package leader_election

import (
	"context"
	"fmt" //"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/appscode/go/ioutil"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/clientcmd"
	_ "github.com/lib/pq"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

const (
	RolePrimary = "primary"
	RoleReplica = "replica"
)

type pgOpCommand int

const (
	startMasterEmpty = iota
	startMasterRecovery
	startSlave
	createRecoveryTrigger
	removeRecoveryTrigger
	promoteToMaster
	promoteToSlave
	raiseError
	raiseFatalError
)

type pgOpErrorType int

const (
	masterUnreachable = iota
	masterNotFunctional
	backupUnreachable
	lostSync
	walgError
	noLeader
)

type pgOpError struct {
	ErrorType pgOpErrorType
	ErrorText string
}

func appendFile(filename string, lines []string) error {
	log.Printf("appendFile: %s", filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	content := strings.Join(lines, "\n")
	if _, err = f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


func runCmd(env []string, cmdname string, params ...string) bool {
	log.Println("runCmd:", cmdname, params)
	cmd := exec.Command(cmdname, params...)

	// set env variables
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)

	log.Printf("runCmd: env %v", env)
	// set stdout, stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// if errors
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return false
	}
	return true
}


func postgresInitDB() {
	// TODO: Search INITDB directory for init files
}

func RunLeaderElection() {

	log.Println("Leader election started")

	leaderElectionLease := 3 * time.Second

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Change owner of Postgres data directory
	if setPermission() == false {
		log.Println("RunLeaderElection: can't chown data directory")
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	parts := strings.Split(hostname, "-")
	statefulSetName := strings.Join(parts[:len(parts)-1], "-")

	log.Printf(`We want "%v" as our leader`, hostname)

	config, err := restclient.InClusterConfig()
	if err != nil {
		log.Println(err)
	}
	clientcmd.Fix(config)

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
	}

	configMap := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetLeaderLockName(statefulSetName),
			Namespace: namespace,
		},
	}
	if _, err := kubeClient.CoreV1().ConfigMaps(namespace).Create(configMap); err != nil && !kerr.IsAlreadyExists(err) {
		log.Println(err)
	}

	resLock := &resourcelock.ConfigMapLock{
		ConfigMapMeta: configMap.ObjectMeta,
		Client:        kubeClient.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      hostname,
			EventRecorder: &record.FakeRecorder{},
		},
	}

	commandsBus := make(chan pgOpCommand)
	recoverySuccessful := make(chan bool)
	runningFirstTime := true
	databaseRestored := false
	mlCtx, mlCancel := context.WithCancel(context.Background())
	defer mlCancel()

	go func() {
		leaderelection.RunOrDie(context.Background(), leaderelection.LeaderElectionConfig{
			Lock:          resLock,
			LeaseDuration: leaderElectionLease,
			RenewDeadline: leaderElectionLease * 2 / 3,
			RetryPeriod:   leaderElectionLease / 3,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					log.Println("Received message to start as master")
				},
				OnStoppedLeading: func() {
					log.Println("Lost leadership, now quit")
					os.Exit(1)
				},
				OnNewLeader: func(identity string) {
					log.Printf("We got new leader - %v!", identity)
					statefulSet, err := kubeClient.AppsV1().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
					if err != nil {
						log.Println(err)
					}

					pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{
						LabelSelector: metav1.FormatLabelSelector(statefulSet.Spec.Selector),
					})
					if err != nil {
						log.Println(err)
					}

					for _, pod := range pods.Items {
						role := RoleReplica
						if pod.Name == identity {
							role = RolePrimary
						}
						_, _, err = core_util.PatchPod(kubeClient, &pod, func(in *core.Pod) *core.Pod {
							in.Labels["kubedb.com/role"] = role
							return in
						})
					}

					if runningFirstTime {
						runningFirstTime = false
						log.Println("Pod started first time")
						if identity == hostname {
							log.Println("Pod started as master server")
							// OS env variable RESTORE contains true or false
							if getEnv("RESTORE", "false") == "true" {
								log.Println("$RESTORE is true, sending signal to start recovery")
								// create context with cancel go master loop go routine
								commandsBus <- startMasterRecovery
								recoveryComplete := <-recoverySuccessful
								if recoveryComplete == true {
									databaseRestored = true
								} else {
									databaseRestored = false
								}

							} else {
								log.Println("$RESTORE is false, sending signal to start empty master")
								commandsBus <- startMasterEmpty
							}
						} else {
							log.Println("Pod started as slave server")
							commandsBus <- startSlave
							recoveryComplete := <-recoverySuccessful
							if recoveryComplete == true {
								databaseRestored = true
							} else {
								databaseRestored = false
							}
						}

					} else {
						mlCancel()
						go masterLoop(mlCtx, commandsBus, recoverySuccessful)
						if databaseRestored == true {
							if identity == hostname {
								log.Println("Creating recovery trigger")
								commandsBus <- createRecoveryTrigger
							}
						} else {
							if identity == hostname {
								log.Println("Database not restored! Cant promote slave to master")
								log.Println("Trying to restore database from backup")
								commandsBus <- startMasterRecovery
							}
						}
					}

				},
			},
		})
	}()

	go masterLoop(mlCtx, commandsBus, recoverySuccessful)
	select {}
}

func masterLoop(ctx context.Context, commandsBus chan pgOpCommand, recoverySuccessful chan bool) {
	log.Println("master loop welcomes you")
	exitLoop := false
	for exitLoop == false {
		select {
		case <-ctx.Done():
			exitLoop = true
			break
		case operatorCommand := <-commandsBus:
			// receive message
			if operatorCommand == startMasterEmpty {
				log.Println("master loop: Received command start as master:", operatorCommand)
				// some actions before start as master
				dataDirectoryCleanup()
				postgresMakeConfigs(RolePrimary)
				postgresMakeEmptyDB()
				setPermission()
				execPostgresAction("start")
			}
			if operatorCommand == startMasterRecovery {
				log.Println("master loop: Received command start recovery:", operatorCommand)
				recoverySuccessful <- restoreMasterFromBackup()

				// if alien backup used set user and password to current deployment credentials
				// be sure that current deployment role have rights to all required resources
				ctxMasterRecovery, cancelMasterRecovery := context.WithCancel(ctx)
				defer cancelMasterRecovery()
				if isPostgresOnline(ctxMasterRecovery, "localhost", true) {
					appendFile("/tmp/pg-failover-trigger", []string{})
					recovery_done_file := getEnv("PGDATA", "/var/pv/data") + "/recovery.done"
					for {
						if _, err := os.Stat(recovery_done_file); !os.IsNotExist(err) {
							break
						}
						log.Println("master loop: Waiting recovery.done to be created")
						time.Sleep(time.Second)
					}
					setPosgresUserPassword(getEnv("POSTGRES_USER", "postgres"), getEnv("POSTGRES_PASSWORD", "postgres"))
				}

			}
			if operatorCommand == startSlave {
				ctxSlave, cancelSlave := context.WithCancel(ctx)
				defer cancelSlave()
				if isPostgresOnline(ctxSlave, getEnv("PRIMARY_HOST", ""), true) {
					log.Println("master loop: Received command start as slave:", operatorCommand)
					// some actions before start as slave
					dataDirectoryCleanup()
					restoreComplete := execBaseBackup()
					postgresMakeConfigs(RoleReplica)
					setPermission()
					if restoreComplete == false {
						recoverySuccessful <- true
					} else {
						recoverySuccessful <- false
					}
					execPostgresAction("start")
				}
			}
			if operatorCommand == createRecoveryTrigger {
				log.Println("master loop: Received command create failover trigger:", operatorCommand)

				if !ioutil.WriteString("/tmp/pg-failover-trigger", "") {
					log.Println("master loop: Failed to create trigger file")
				}
			}

		default:
			// nothing, sleep to lower cpu usage
			time.Sleep(time.Second)
		}
	}
}

func setPermission() bool {
	log.Println("setPermission: chown data directory")
	var env []string
	runCmd(env, "chown", "-R", "postgres:postgres", getEnv("PGDATA", "/var/pv/data"))
	runCmd(env, "chmod", "-R", "700", getEnv("PGDATA", "/var/pv/data"))
	return true
}

func GetLeaderLockName(offshootName string) string {
	return fmt.Sprintf("%s-leader-lock", offshootName)
}
