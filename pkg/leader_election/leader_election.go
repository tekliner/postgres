package leader_election

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pgConnString() string {

	hostname := getEnv("PRIMARY_HOST", "localhost")
	username := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")

	info := fmt.Sprintf("host=%s port=%d dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		hostname,
		5432,
		"postgres",
		"disable",
		username,
		password,
	)
	log.Printf("posgres connection string: %v", info)
	return info
}

func isPgMasterOnline() bool {
	log.Println("Checking connection to master")

	if db, err := sql.Open("postgres", pgConnString()); db != nil {
		defer db.Close()
		if _, err = db.Exec("SELECT 1;"); err == nil {
			return true
		}
		log.Println("Checking connection to master: query error")
	} else {
		log.Println("Checking connection to master: connection error")
	}

	time.Sleep(time.Second * 5)
	return false
}

func waitPgMasterOnline() {
	// authung! dangerous function, can't stop if master offline
	// use with go and channels
	for {
		log.Println("Waiting connection to master")

		if db, err := sql.Open("postgres", pgConnString()); db != nil {
			defer db.Close()
			if _, err = db.Exec("SELECT 1;"); err == nil {
				break
			}
			log.Println("Checking connection to master: query error")
		} else {
			log.Println("Checking connection to master: connection error")
		}

		time.Sleep(time.Second * 5)

	}

}

func runCmd(env []string, cmdname string, params ...string) {
	log.Println("runCmd:", cmdname, params)
	cmd := exec.Command(cmdname, params...)

	// set env variables
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)

	// set stdout, stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// if errors
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	os.Exit(1)
}

func execWalgAction(walgCommand string, params ...string) {
	var env []string
	env = append(env, fmt.Sprintf("WALE_S3_PREFIX=%s", getEnv("ARCHIVE_S3_PREFIX", "")))
	// auth for wal-g
	env = append(env, fmt.Sprintf("PGUSER=%s", getEnv("POSTGRES_USER", "")))
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", getEnv("POSTGRES_PASSWORD", "")))

	awsKeyFile := "/srv/wal-g/archive/secrets/AWS_ACCESS_KEY_ID"
	awsSecretFile := "/srv/wal-g/archive/secrets/AWS_SECRET_ACCESS_KEY"

	awsKey, err := ioutil.ReadFile(awsKeyFile)
	// aws key file ansent
	if err != nil {
		log.Println(err)
	}

	awsSecret, err := ioutil.ReadFile(awsSecretFile)
	// aws secret file absent
	if err != nil {
		log.Println(err)
	}
	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", awsKey))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", awsSecret))
	// need to forward "wal-g", walgCommand, params...
	arg := []string{"postgres", "wal-g"}
	arg = append(arg, walgCommand)
	arg = append(arg, params...)
	runCmd(env, "su-exec", arg...)
}

func execPostgresAction(action string) {
	var env []string
	runCmd(env, "su-exec", "postgres", "pg_ctl", "-D", getEnv("PGDaTA", "/var/pv/data"), "-w", action)
}

// TODO: generate PG configs

func RunLeaderElection() {

	log.Println("Leader election started")

	leaderElectionLease := 3 * time.Second

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Change owner of Postgres data directory
	if err := setPermission(); err != nil {
		log.Fatalln(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	parts := strings.Split(hostname, "-")
	statefulSetName := strings.Join(parts[:len(parts)-1], "-")

	log.Printf(`We want "%v" as our leader`, hostname)

	config, err := restclient.InClusterConfig()
	if err != nil {
		log.Fatalln(err)
	}
	clientcmd.Fix(config)

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	configMap := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetLeaderLockName(statefulSetName),
			Namespace: namespace,
		},
	}
	if _, err := kubeClient.CoreV1().ConfigMaps(namespace).Create(configMap); err != nil && !kerr.IsAlreadyExists(err) {
		log.Fatalln(err)
	}

	resLock := &resourcelock.ConfigMapLock{
		ConfigMapMeta: configMap.ObjectMeta,
		Client:        kubeClient.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      hostname,
			EventRecorder: &record.FakeRecorder{},
		},
	}

	runningFirstTime := true

	type pgOpCommand int
	const (
		startMaster           pgOpCommand = 0
		startSlave            pgOpCommand = 1
		startRecovery         pgOpCommand = 2
		createRecoveryTrigger pgOpCommand = 3
		removeRecoveryTrigger pgOpCommand = 4
		promoteToMaster       pgOpCommand = 5
		promoteToSlave        pgOpCommand = 6
		raiseError            pgOpCommand = 7
		raiseFatalError       pgOpCommand = 8
	)

	type pgOpErrorType int
	const (
		masterUnreachable   pgOpErrorType = 0
		masterNotFunctional pgOpErrorType = 1
		backupUnreachable   pgOpErrorType = 2
		lostSync            pgOpErrorType = 3
		walgError           pgOpErrorType = 4
		noLeader            pgOpErrorType = 5
	)

	type pgOpError struct {
		ErrorType pgOpErrorType
		ErrorText string
	}

	messagesBus := make(chan string)
	commandsBus := make(chan pgOpCommand)
	errorsBus := make(chan pgOpError)

	go func() {
		// master loop
		for {
			select {
			case operatorCommand := <-commandsBus:
				// receive message
				if operatorCommand == startMaster {
					log.Println("Received command start as master:", operatorCommand)
					// some actions before start as master
					// TODO: run config generation
					execPostgresAction("start")
				}
				if operatorCommand == startSlave {
					log.Println("Received command start as slave:", operatorCommand)
					// some actions before start as slave
					// TODO: run config generation
					execPostgresAction("start")
				}
				if operatorCommand == startRecovery {
					log.Println("Received command start recovery:", operatorCommand)
					// some actions before start recovery
					execPostgresAction("start")
					// some actions to start recovery

					// wait backup list to be fetched
					// env:
					// PITR=${PITR:-false}
					// TARGET_INCLUSIVE=${TARGET_INCLUSIVE:-true}
					// TARGET_TIME=${TARGET_TIME:-}
					// TARGET_TIMELINE=${TARGET_TIMELINE:-}
					// TARGET_XID=${TARGET_XID:-}
					execWalgAction("backup-list")
					execWalgAction("backup-fetch", getEnv("PGDaTA", "/var/pv/data"), getEnv("BACKUP_NAME", "LATEST"))
					execPostgresAction("stop")
				}
			case operatorerror := <-errorsBus:
				//
			default:
				// nothing
				log.Println("Can't connect to master server")
			}
		}
	}()

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
					// TODO: change hardcoded directory in setPermission to PGDATA

					log.Printf("We got new leader - %v!", identity)
					statefulSet, err := kubeClient.AppsV1().StatefulSets(namespace).Get(statefulSetName, metav1.GetOptions{})
					if err != nil {
						log.Fatalln(err)
					}

					pods, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{
						LabelSelector: metav1.FormatLabelSelector(statefulSet.Spec.Selector),
					})
					if err != nil {
						log.Fatalln(err)
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

					log.Println("Set default role to replica")

					role := RoleReplica
					if identity == hostname {
						role = RolePrimary
					}

					for role == RoleReplica {
						log.Println("Waiting connection to master")

						if db, err := sql.Open("postgres", pgConnString()); db != nil {
							defer db.Close()
							if _, err = db.Exec("SELECT 1;"); err == nil {
								break
							}
							log.Println("Waiting connection to master: query error")
						} else {
							log.Println("Waiting connection to master: connection error")
						}

						time.Sleep(time.Second * 5)

						select {
						case trigger := <-startAsMaster:
							log.Println("Got leadership:", trigger)
							role = RolePrimary
						default:
							// nothing
							log.Println("Can't connect to master server")
						}
					}

					if runningFirstTime {
						runningFirstTime = false
						go func() {
							log.Println("Starting script:", role)
							// su-exec postgres /scripts/primary/run.sh
							cmd := exec.Command("su-exec", "postgres", fmt.Sprintf("/scripts/%s/run.sh", role))
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr

							if err = cmd.Run(); err != nil {
								log.Println(err)
							}
							os.Exit(1)
						}()
					} else {
						if identity == hostname {
							log.Fatalln("Creating trigger file")
							if !ioutil.WriteString("/tmp/pg-failover-trigger", "") {
								log.Fatalln("Failed to create trigger file")
							}
						}
					}
				},
			},
		})
	}()

	select {}
}

func setPermission() error {
	u, err := user.Lookup("postgres")
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}
	err = os.Chown("/var/pv", uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func GetLeaderLockName(offshootName string) string {
	return fmt.Sprintf("%s-leader-lock", offshootName)
}
