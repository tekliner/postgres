package report

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	pg "github.com/lib/pq"
	api "github.com/tekliner/apimachinery/apis/kubedb/v1alpha1"
)

func newXormEngine(username, password, host string, port int, dbName string) (*xorm.Engine, error) {
	cnnstr := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		username, password, host, port, dbName)

	engine, err := xorm.NewEngine("postgres", cnnstr)
	if err != nil {
		return nil, err
	}

	engine.SetMaxIdleConns(0)
	engine.DB().SetConnMaxLifetime(10 * time.Minute)
	engine.ShowSQL(false)
	engine.Logger().SetLevel(core.LOG_ERR)
	return engine, nil
}

func getAllDatabase(engine *xorm.Engine) ([]string, error) {
	defer engine.Close()
	engine.ShowSQL(true)
	session := engine.NewSession()
	defer session.Close()

	rows, err := session.Query("SELECT datname FROM pg_database where datistemplate = false")
	if err != nil {
		return nil, err
	}

	databases := make([]string, 0)

	for _, row := range rows {
		databases = append(databases, string(row["datname"]))
	}
	return databases, nil
}

func getDataFromDB(engine *xorm.Engine) (*api.PostgresSummary, error) {
	defer engine.Close()
	engine.ShowSQL(true)
	session := engine.NewSession()
	defer session.Close()

	schemaRowSlice, err := session.Query("select schema_name from information_schema.schemata")
	if err != nil {
		return nil, err
	}

	schemaList := make(map[string]*api.PostgresSchemaInfo, 0)
	for _, row := range schemaRowSlice {
		schemaName := string(row["schema_name"])
		schemaInfo, err := getDataFromSchema(session, schemaName)
		if err != nil {
			return nil, err
		}
		schemaList[schemaName] = schemaInfo
	}

	return &api.PostgresSummary{
		Schema: schemaList,
	}, nil
}

func getDataFromSchema(session *xorm.Session, schemaName string) (*api.PostgresSchemaInfo, error) {
	tableRowSlice, err := session.Query("SELECT tablename FROM pg_tables where schemaname=$1", schemaName)
	if err != nil {
		log.Fatalln(err)
	}

	schemaInfo := &api.PostgresSchemaInfo{
		Table: make(map[string]*api.PostgresTableInfo),
	}

	for _, row := range tableRowSlice {
		for _, val := range row {
			tableName := string(val)
			tableInfo, err := getDataFromTable(session, schemaName, tableName)
			if err != nil {
				return nil, err
			}
			schemaInfo.Table[tableName] = tableInfo
		}
	}

	return schemaInfo, nil
}

const (
	errorUndefinedColumn  = "undefined_column"
	errorDatatypeMismatch = "datatype_mismatch"
	invalidData           = -1
	TotalRow              = "total_row"
	MaxID                 = "max_id"
	NextID                = "next_id"
)

func getDataFromTable(session *xorm.Session, schemaName, tableName string) (*api.PostgresTableInfo, error) {
	table := fmt.Sprintf(`"%v".%v`, schemaName, tableName)
	dataRows, err := session.Query(fmt.Sprintf(`SELECT count(*) as total_row, coalesce(max(id),0) as max_id FROM %v`, table))

	var totalRow, maxID, nextID int64
	var errorName string

	if driverErr, ok := err.(*pg.Error); ok {
		errorName = driverErr.Code.Name()
		if errorName == errorUndefinedColumn || errorName == errorDatatypeMismatch {
			dataRows, err = session.Query(fmt.Sprintf("SELECT count(*) as total_row FROM %v", table))
			if err != nil {
				return &api.PostgresTableInfo{}, err
			}

			if totalRow, err = strconv.ParseInt(string(dataRows[0][TotalRow]), 10, 64); err != nil {
				return &api.PostgresTableInfo{}, err
			}
			maxID = invalidData
			nextID = invalidData

		} else {
			return &api.PostgresTableInfo{}, err
		}
	} else {
		if len(dataRows) == 0 {
			log.Println("Data missing: ", err)
			totalRow = invalidData
			maxID = invalidData
			nextID = invalidData
		} else {
			if totalRow, err = strconv.ParseInt(string(dataRows[0][TotalRow]), 10, 64); err != nil {
				return &api.PostgresTableInfo{}, err
			}

			if maxID, err = strconv.ParseInt(string(dataRows[0][MaxID]), 10, 64); err != nil {
				return &api.PostgresTableInfo{}, err
			}

			dataRows, err = session.Query(fmt.Sprintf(`select (last_value+1) as next_id from %v_id_seq`, table))
			if err != nil {
				return &api.PostgresTableInfo{}, err
			}
			if len(dataRows) == 0 {
				nextID = invalidData
			} else {
				if nextID, err = strconv.ParseInt(string(dataRows[0][NextID]), 10, 64); err != nil {
					return &api.PostgresTableInfo{}, err
				}
			}
		}
	}

	return &api.PostgresTableInfo{
		TotalRow: totalRow,
		MaxID:    maxID,
		NextID:   nextID,
	}, nil
}
