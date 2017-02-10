package mig // import "iris.arke.works/forum/db/mig"

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASS")
	viper.BindEnv("POSTGRES_ABORT_MIG")
}

func TestDB(t *testing.T) {
	assert := assert.New(t)

	var connString string

	if !viper.IsSet("POSTGRES_HOST") {
		t.Log("DB not set, aborting Database Test")
		return
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
	db, err := sql.Open("postgres", connString)
	assert.NoError(err)
	if err != nil {
		return
	}
	if err := db.Ping(); err != nil {
		t.Log("Could not ping DB, aborting test silently")
		t.Log(err)
		return
	}
	defer db.Close()

	if viper.Get("POSTGRES_ABORT_MIG") == "YES" {
		t.Log("Local testing aborted")
		return
	}

	migDB := OpenFromPGConn(db)
	assert.NoError(migDB.CheckAndLoadTables())
	// Check again to see if existing tables don't produce errors
	assert.NoError(migDB.CheckAndLoadTables())

	graph := NewGraph()
	assert.NoError(graph.Load("arke"))

	graph, err = graph.GetTargetSubgraph("default")
	assert.NoError(err)

	executedUnits, err := migDB.GetExecutedUnits()
	assert.NoError(err)
	assert.NoError(graph.MarkNodesRun(executedUnits...))

	for nodes := graph.GetAllRunnableNodes(); len(nodes) > 0; nodes = graph.GetAllRunnableNodes() {
		for _, v := range nodes {
			node, err := graph.GetUnit(v)
			assert.NoError(err)
			_, err = db.Exec(node.SQL.Postgres)
			assert.NoError(err)
			assert.NoError(migDB.MarkExecuted(node))
		}
		assert.NoError(graph.MarkNodesRun(nodes...))
		assert.False(graph.IsStuck())
	}
}
