package leader_election

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

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
	log.Printf("postgres connection string: %v", info)
	return info
}

func setPosgresUserPassword(username, password string) {
	log.Printf("Trying to set password to Postgres user: %s", username)

	if db, err := sql.Open("postgres", pgConnString("localhost")); db != nil {
		defer db.Close()
		sqlQuery := fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s';", username, password)
		if _, err = db.Exec(sqlQuery); err == nil {
			log.Printf("Password successfully set to %s", password)
		}
		log.Println("query error")
	} else {
		log.Println("connection error")
	}
}

func isPostgresOnline(ctx context.Context, hostname string, wait bool) bool {
	// authung! dangerous function
	//if wait == true function will wait until connection established
	returnValue := false
	exitLoop := false
	for exitLoop == false {
		select {
		case <-ctx.Done():
			exitLoop = true
			returnValue = false
			break
		}

		log.Println("Checking connection to master")

		if db, err := sql.Open("postgres", pgConnString(hostname)); db != nil {
			defer db.Close()
			if _, err = db.Exec("SELECT 1;"); err == nil {
				returnValue = true
				if wait {
					break
				}
			}
			log.Println("query error")
			db.Close()
		} else {
			log.Println("connection error")
		}
		if !wait {
			break
		}
		time.Sleep(time.Second * 60)
	}
	return returnValue
}

func dataDirectoryCleanup() {
	log.Println("dataDirectoryCleanup: Removing and creating data directory")
	PGDATA := getEnv("PGDATA", "/var/pv/data")
	os.RemoveAll(PGDATA)
	os.MkdirAll(PGDATA, 0755)
	setPermission()
}

func dataDirectoryCreateAfterWalg() {
	log.Println("dataDirectoryCreateAfterWalg: creating directories after wal-g")
	PGDATA := getEnv("PGDATA", "/var/pv/data")
	dirList := []string{
		"/pg_tblspc",
		"/pg_twophase",
		"/pg_stat",
		"/pg_commit_ts",
		"/pg_logical/snapshots",
		"/pg_logical/mappings",
	}
	for _, dir := range dirList {
		os.MkdirAll(PGDATA+dir, 0700)
	}
}

func execWalgAction(walgCommand string, params ...string) bool {
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
	arg := []string{"root", "wal-g"}
	arg = append(arg, walgCommand)
	arg = append(arg, params...)
	return runCmd(env, "su-exec", arg...)
}

func execBaseBackup() bool {
	log.Println("execBaseBackup: running pg_basebackup")
	var env []string
	env = append(env, fmt.Sprintf("PGUSER=%s", getEnv("POSTGRES_USER", "postgres")))
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", getEnv("POSTGRES_PASSWORD", "postgres")))
	pgdata := getEnv("PGDaTA", "/var/pv/data")
	pghost := fmt.Sprintf("--host=%s", getEnv("PRIMARY_HOST", ""))
	return runCmd(env, "pg_basebackup", "-X", "fetch", "--no-password", "--pgdata", pgdata, pghost)
}

func execPostgresAction(action string) {
	log.Printf("execPostgresAction: %s", action)
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
	runCmd(env, "su-exec", "postgres", "pg_ctl", "-D", getEnv("PGDATA", "/var/pv/data"), "-w", action)
}

func postgresMakeEmptyDB() {
	log.Println("postgresMakeEmptyDB: Create empty database for postgres")
	var env []string
	runCmd(env, "initdb", fmt.Sprintf("--pgdata=%s", getEnv("PGDATA", "/var/pv/data")))
}

func postgresMakeConfigs(role string) {
	log.Printf("Create config files for postgres, role: %s", role)
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

func restoreMasterFromBackup() bool {
	// absolutely clean data directory
	dataDirectoryCleanup()
	// some actions to start recovery
	execWalgAction("backup-list")
	restoreComplete := execWalgAction("backup-fetch", getEnv("PGDATA", "/var/pv/data"), getEnv("BACKUP_NAME", "LATEST"))
	dataDirectoryCreateAfterWalg()
	postgresMakeConfigs(RolePrimary)
	// additional lines to recovery.conf
	lines := []string{}
	if getEnv("PITR", "") != "" {

		target_inclusive := getEnv("TARGET_INCLUSIVE", "true")
		target_time := getEnv("TARGET_TIME", "")
		target_timeline := getEnv("TARGET_TIMELINE", "")
		target_xid := getEnv("TARGET_XID", "")

		lines = []string{
			fmt.Sprintf("recovery_target_inclusive = '%s'", target_inclusive),
			"recovery_target_action = 'promote'",
		}
		if target_time != "" {
			lines = append(lines, fmt.Sprintf("recovery_target_time = '%s'", target_time))
		}
		if target_timeline != "" {
			lines = append(lines, fmt.Sprintf("recovery_target_timeline = '%s'", target_timeline))
		}
		if target_xid != "" {
			lines = append(lines, fmt.Sprintf("recovery_target_xid = '%s'", target_xid))
		}

	}
	lines = append(lines, "restore_command = 'wal-g wal-fetch %f %p'")
	appendFile(getEnv("PGDATA", "/var/pv/data")+"/recovery.conf", lines)
	os.Remove(getEnv("PGDATA", "/var/pv/data") + "/recovery.done")
	setPermission()
	go execPostgresAction("start")
	// backup done, start Postgres
	if restoreComplete == false {
		return true
	}

	return false
}
