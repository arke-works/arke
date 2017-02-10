package models

import (
	"database/sql"
	"flag"
	"fmt"
	"iris.arke.works/forum/db/mig"
	"iris.arke.works/forum/snowflakes"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"os"
	"testing"
	"time"
)

func init() {
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASS")
	viper.BindEnv("POSTGRES_ABORT_MIG")
}

var db *sql.DB

var snow *snowflakes.Generator

func TestMain(m *testing.M) {

	flag.Parse()

	var connString string

	snow = &snowflakes.Generator{
		StartTime:  time.Date(1998, time.November, 19, 0, 0, 0, 0, time.UTC).Unix(),
		InstanceID: 1,
	}

	if !viper.IsSet("POSTGRES_HOST") {
		println("PostgreSQL Host not setup")
		os.Exit(1)
	}

	if viper.IsSet("POSTGRES_PASS") && len(viper.GetString("POSTGRES_PASS")) > 0 {
		connString = fmt.Sprintf(
			"postgres://%s:%s@%s/?sslmode=disable",
			viper.Get("POSTGRES_USER"),
			viper.Get("POSTGRES_PASS"),
			viper.Get("POSTGRES_HOST"),
		)
	} else {
		connString = fmt.Sprintf(
			"postgres://%s@%s/?sslmode=disable",
			viper.Get("POSTGRES_USER"),
			viper.Get("POSTGRES_HOST"),
		)
	}
	var err error
	db, err = sql.Open("postgres", connString)
	if err != nil {
		println("DB Connection failed", connString)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		println("DB Ping failed")
		os.Exit(1)
	}
	defer db.Close()

	migDB := mig.OpenFromPGConn(db)
	err = migDB.CheckAndLoadTables()
	if err != nil {
		os.Exit(1)
	}

	executedUnits, err := migDB.GetExecutedUnits()

	graph := mig.NewGraph()
	graph.Load("mig/arke")
	graph.MarkNodesRun(executedUnits...)

	for nodes := graph.GetAllRunnableNodes(); len(nodes) > 0; nodes = graph.GetAllRunnableNodes() {
		for _, v := range nodes {
			node, _ := graph.GetUnit(v)
			_, _ = db.Exec(node.SQL.Postgres)
			migDB.MarkExecuted(node)
		}
		graph.MarkNodesRun(nodes...)
		graph.IsStuck()
	}

	boil.SetDB(db)

	os.Exit(m.Run())
}
