package leader_election

import (
	"context"
	"database/sql"
	"fmt"
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

//Юзать коллбэки на начало и конец лидерства, а тот, где простыня не юзать.
//Перенести чек на доступность мастера в го и мочь отменить его через context.

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pg_conn_string() string {

	hostname := getEnv("PRIMARY_HOST", "localhost")
	username := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		hostname,
		5432,
		"postgres",
		"disable",
		username,
		password,
	)
	return info
}

func RunLeaderElection() {

	fmt.Println("Leader election started")

	leaderElectionLease := 3 * time.Second

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Change owner of Postgres data directory
	// нахуя это делать в коде, если конфик собирается скриптом?
	if err := setPermission(); err != nil {
		log.Fatalln(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	parts := strings.Split(hostname, "-")
	statefulSetName := strings.Join(parts[:len(parts)-1], "-")

	fmt.Println(fmt.Sprintf(`We want "%v" as our leader`, hostname))

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
	start_as_master := make(chan int, 10)

	go func() {
		leaderelection.RunOrDie(context.Background(), leaderelection.LeaderElectionConfig{
			Lock:          resLock,
			LeaseDuration: leaderElectionLease,
			RenewDeadline: leaderElectionLease * 2 / 3,
			RetryPeriod:   leaderElectionLease / 3,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					fmt.Println("Received message to start as master")
					start_as_master <- 1
				},
				OnStoppedLeading: func() {
					fmt.Println("Lost leadership, now quit")
					os.Exit(1)
				},
				// вызывается при поиске лидера, если его нет попытаться его получить
				OnNewLeader: func(identity string) {
					fmt.Println("We got new leader!")
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

					fmt.Println("Set default role to replica")

					role := RoleReplica
					if identity == hostname {
						role = RolePrimary
					}

					// добавить тест мастера с помощью select 1
					// при ошибке опроса через канал сообщить о том, что мастер труп в другие рутины
					// начать восстановление как мастера в рутине OnStartedLeading

					for {
						query_error := false
						conn_error := false

						fmt.Println("Checking connection to master")

						if conn_error || query_error {

							db, err := sql.Open("postgres", pg_conn_string())
							if err != nil {
								conn_error = true
								fmt.Println("Checking connection to master: connection error")
							}

							defer db.Close()

							_, err = db.Exec("SELECT 1;")
							if err != nil {
								query_error = true
								fmt.Println("Checking connection to master: query error")
							}

							time.Sleep(time.Second * 5)

							select {
							case trigger := <-start_as_master:
								fmt.Println("Got leadership:", trigger)
								role = RolePrimary
							default:
								// nothing
								fmt.Println("Can't connect to master server")
							}
						}
					}

					if runningFirstTime {
						runningFirstTime = false
						go func() {
							fmt.Println("Starting script:", role)
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
							if !ioutil.WriteString("/tmp/pg-failover-trigger", "") {
								log.Fatalln("Failed to create trigger file")
								fmt.Println("Failed to create trigger file")
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
