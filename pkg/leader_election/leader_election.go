package leader_election

import (
	"context"
	"database/sql"
	"fmt" //"io/ioutil"
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

func appendFile(filename string, lines []string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
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

func dataDirectoryCleanup() {
	PGDATA := getEnv("PGDATA", "/var/pv/data")
	os.RemoveAll(PGDATA)
	os.MkdirAll(PGDATA, 0700)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pgConnString(hostname string) string {

	//hostname := getEnv("PRIMARY_HOST", "localhost")
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
	log.Printf("pgConnString: posgres connection string: %v", info)
	return info
}

func setPosgresUserPassword(username, password string) {
	log.Printf("setPostgresUserPassword: Trying to set password to Postgres user: %s", username)

	if db, err := sql.Open("postgres", pgConnString("localhost")); db != nil {
		defer db.Close()
		sqlQuery := fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s';", username, password)
		if _, err = db.Exec(sqlQuery); err == nil {
			log.Printf("setPostgresUserPassword: Password successfully set to %s", password)
		}
		log.Println("setPostgresUserPassword: query error")
	} else {
		log.Println("setPostgresUserPassword: connection error")
	}
}

func isPostgresOnline(hostname string, wait bool) bool {
	// authung! dangerous function
	//if wait == true function will wait until connection established
	returnValue := false
	for {
		log.Println("isPgMasterOnline: Checking connection to master")

		if db, err := sql.Open("postgres", pgConnString(hostname)); db != nil {
			defer db.Close()
			if _, err = db.Exec("SELECT 1;"); err == nil {
				returnValue = true
				if wait {
					break
				}
			}
			log.Println("isPgMasterOnline: query error")
		} else {
			log.Println("isPgMasterOnline: connection error")
		}
		if !wait {
			break
		}
		time.Sleep(time.Second * 5)
	}
	return returnValue
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

func execBaseBackup() {
	var env []string
	pgdata := getEnv("PGDaTA", "/var/pv/data")
	pguser := fmt.Sprintf("--username=%s", getEnv("POSTGRES_USER", "postgres"))
	pghost := fmt.Sprintf("--host=%s", getEnv("PRIMARY_HOST", ""))
	runCmd(env, "pg_basebackup", "-X", "fetch", "--no-password", "--pgdata", pgdata, pguser, pghost)
}

func execPostgresAction(action string) {
	var env []string
	runCmd(env, "su-exec", "postgres", "pg_ctl", "-D", getEnv("PGDaTA", "/var/pv/data"), "-w", action)
}

func postgresMakeEmptyDB() {
	// Create empty database for postgres after dir cleanup and before start
	var env []string
	runCmd(env, "initdb", fmt.Sprintf("--pgdata=%s", getEnv("PGDATA", "/var/pv/data")))
}

func postgresMakeConfigs(role string) {
	// Create config files for postgres after dir cleanup and before start
	if role == RolePrimary {
		var env []string
		// copy template to /tmp
		runCmd(env, "cp", "/scripts/primary/postgresql.conf", "/tmp/")

		// append config
		lines := []string{
			"wal_level = replica",
			"max_wal_senders = 99",
			"wal_keep_segments = 32",
		}
		appendFile("/tmp/postgresql.conf", lines)

		// move configs to PGDATA
		runCmd(env, "mv", "/tmp/postgresql.conf", getEnv("PGDATA", "/var/pv/data"))
		runCmd(env, "mv", "/scripts/primary/pg_hba.conf", getEnv("PGDATA", "/var/pv/data"))
	}
	if role == RoleReplica {
		var env []string
		// copy template to /tmp
		runCmd(env, "cp", "/scripts/replica/recovery.conf", "/tmp/")

		// append recovery.conf
		lines := []string{
			"recovery_target_timeline = 'latest'",
			fmt.Sprintf("archive_cleanup_command = 'pg_archivecleanup %s %r'", getEnv("PGWAL", "")),
			fmt.Sprintf("primary_conninfo = 'application_name=%s host=%s'", getEnv("HOSTNAME", ""), getEnv("PRIMARY_HOST", "")),
		}
		appendFile("/tmp/recovery.conf", lines)

		// append postgresql.conf
		runCmd(env, "cp", "/scripts/primary/postgresql.conf", "/tmp/")
		lines = []string{
			"wal_level = replica",
			"max_wal_senders = 99",
			"wal_keep_segments = 32",
		}
		if getEnv("STANDBY", "") == "hot" {
			lines = append(lines, "hot_standby = on")
		}
		if getEnv("STREAMING", "") == "synchronous" {
			// setup synchronous streaming replication
			lines = append(lines, "synchronous_commit = remote_write")
			lines = append(lines, "synchronous_standby_names = '*'")
		}
		if getEnv("ARCHIVE", "") == "wal-g" {
			lines = append(lines, "archive_command = 'wal-g wal-push %p'")
			lines = append(lines, "archive_timeout = 60")
			lines = append(lines, "archive_mode = always")
		}
		appendFile("/tmp/postgresql.conf", lines)

		// move configs to PGDATA
		runCmd(env, "mv", "/tmp/postgresql.conf", getEnv("PGDATA", "/var/pv/data"))
		runCmd(env, "mv", "/tmp/recovery.conf", getEnv("PGDATA", "/var/pv/data"))
		runCmd(env, "mv", "/scripts/primary/pg_hba.conf", getEnv("PGDATA", "/var/pv/data"))
	}
}

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

	type pgOpCommand int
	const (
		startMasterEmpty      pgOpCommand = 0
		startMasterRecovery   pgOpCommand = 1
		startSlave            pgOpCommand = 2
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

	//messagesBus := make(chan string)
	commandsBus := make(chan pgOpCommand)
	errorsBus := make(chan pgOpError)
	runningFirstTime := true

	go func() {
		// master loop
		for {
			select {
			case operatorCommand := <-commandsBus:
				// receive message
				if operatorCommand == startMasterEmpty {
					log.Println("master loop: Received command start as master:", operatorCommand)
					// some actions before start as master
					dataDirectoryCleanup()
					postgresMakeConfigs(RolePrimary)
					postgresMakeEmptyDB()
					execPostgresAction("start")
				}
				if operatorCommand == startMasterRecovery {
					log.Println("master loop: Received command start recovery:", operatorCommand)
					// some actions before start recovery
					dataDirectoryCleanup()
					postgresMakeConfigs(RolePrimary)
					// some actions to start recovery
					execWalgAction("backup-list")
					execWalgAction("backup-fetch", getEnv("PGDaTA", "/var/pv/data"), getEnv("BACKUP_NAME", "LATEST"))
					setPermission()
					// backup done, start Postgres
					go execPostgresAction("start")
					// if alien backup used set user and password to current deployment credentials
					// be sure that current deployment role have rights to all required resources
					if isPostgresOnline("localhost", true) {
						setPosgresUserPassword(getEnv("POSTGRES_USER", "postgres"), getEnv("POSTGRES_PASSWORD", "postgres"))
					}

				}
				if operatorCommand == startSlave {
					log.Println("master loop: Received command start as slave:", operatorCommand)
					// some actions before start as slave
					dataDirectoryCleanup()
					execBaseBackup()
					postgresMakeConfigs(RoleReplica)
					setPermission()
					execPostgresAction("start")
				}
				if operatorCommand == createRecoveryTrigger {
					log.Println("master loop: Received command create failover trigger:", operatorCommand)

					if !ioutil.WriteString("/tmp/pg-failover-trigger", "") {
						log.Fatalln("master loop: Failed to create trigger file")
					}
				}
			case operatorError := <-errorsBus:
				// TODO: create exceptions processing
				log.Printf("master loop: Got error %v", operatorError)

			default:
				// nothing, sleep to lower cpu usage
				time.Sleep(time.Second)
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
					log.Println("RunOrDie: Received message to start as master")
				},
				OnStoppedLeading: func() {
					log.Println("RunOrDie: Lost leadership, now quit")
					os.Exit(1)
				},
				OnNewLeader: func(identity string) {
					log.Printf("RunOrDie: We got new leader - %v!", identity)
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

					if runningFirstTime {
						runningFirstTime = false
						log.Println("RunOrDie: Pod started first time")
						if identity == hostname {
							log.Println("RunOrDie: Pod started as master server")
							// OS env variable RESTORE contains true or false
							if getEnv("RESTORE", "false") == "true" {
								log.Println("RunOrDie: $RESTORE is true, sending signal to start recovery")
								commandsBus <- startMasterRecovery
							} else {
								log.Println("RunOrDie: $RESTORE is false, sending signal to start empty master")
								commandsBus <- startMasterEmpty
							}
						} else {
							log.Println("RunOrDie: Pod started as slave server")
							commandsBus <- startSlave
						}

					} else {
						log.Println("RunOrDie: Pod NOT started first time")
						if identity == hostname {
							commandsBus <- createRecoveryTrigger
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
