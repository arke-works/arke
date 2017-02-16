package cmd // import "iris.arke.works/forum/cmd"

import (
	"database/sql"
	"fmt"
	// Import lib/pq for postgres support
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"iris.arke.works/forum/db/mig"
	"sync"
)

var vapeCmd = &cobra.Command{
	Use:   "vape",
	Short: "Verify and Prepare Environment",
	Long:  "Verifies the configuration, connects to the various end-points and performs various preperation tasks like Database Setup.",
	Run:   run,
}

func init() {
	vapeCmd.Flags().String("migtarget", "default", "Migration Target")
	RootCmd.AddCommand(vapeCmd)
}

func run(cmd *cobra.Command, args []string) {
	log, err := zap.NewProduction()
	if err != nil {
		println("Error while creating logger:", err)
		return
	}
	dbconf := viper.Sub("db.postgres")
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		dbconf.GetString("user"),
		dbconf.GetString("pass"),
		dbconf.GetString("host"),
		dbconf.GetString("dbname"),
		dbconf.GetString("sslmode"))
	log.Info("Opening Database")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error while opening connection to database", zap.Error(err))
		return
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error while verifying connection to database", zap.Error(err))
		return
	}
	defer db.Close()

	log.Info("Verifying and Loading Migration Data from Database")
	migDB := mig.OpenFromPGConn(db)
	err = migDB.CheckAndLoadTables()
	if err != nil {
		log.Fatal("Error while loading migration tables", zap.Error(err))
		return
	}

	log.Info("Loading Migration Units")
	rootGraph := mig.NewGraph()
	err = rootGraph.Load("db/mig/arke")
	if err != nil {
		log.Fatal("Could not load migration data", zap.Error(err))
		return
	}

	log.Info("Validating Migration Units")
	err = rootGraph.ValidateNodes()
	if err != nil {
		log.Fatal("Migration Graph Validation Failed", zap.Error(err))
		return
	}

	target, err := cmd.Flags().GetString("migtarget")
	if err != nil {
		log.Fatal("Migration Target not specified", zap.Error(err))
	}
	log.Info("Loading Migration Target", zap.String("target", target))
	migGraph, err := rootGraph.GetTargetSubgraph(target)
	if err != nil {
		log.Fatal("Could not load Subgraph", zap.Error(err))
		return
	}

	log.Info("Loading already executed Units")
	executedUnits, err := migDB.GetExecutedUnits()
	if err != nil {
		log.Fatal("Error loading executed units", zap.Error(err))
	}

	log.Info("Marking units as already executed", zap.Int("unit_num", len(executedUnits)))
	migGraph.MarkNodesRun(executedUnits...)

	log.Info("Entering Database Transaction")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Transaction failed", zap.Error(err))
		return
	}

	log.Info("Starting Migration", zap.Int("unit_num", migGraph.RemainingSize()))
	hasErrored := false
	for nodes := migGraph.GetAllRunnableNodes(); len(nodes) > 0 && !hasErrored; nodes = migGraph.GetAllRunnableNodes() {
		log.Info("Executing Round", zap.Int("unit_num", len(nodes)), zap.Strings("nodes", nodes))
		wg := sync.WaitGroup{}
		wg.Add(len(nodes))
		for _, v := range nodes {
			go func(v string) {
				defer wg.Done()
				node, err := migGraph.GetUnit(v)
				if err != nil {
					tx.Rollback()
					log.Error("Attempted to migrate non-existant unit", zap.String("unit", v))
					hasErrored = true
					return
				}
				_, err = tx.Exec(node.SQL.Postgres)
				if err != nil {
					tx.Rollback()
					log.Error("Unit failed", zap.String("unit", v), zap.Error(err))
					hasErrored = true
					return
				}
				err = migDB.MarkExecuted(node)
				if err != nil {
					tx.Rollback()
					log.Error("Unit could not be marked as finished", zap.String("unit", v), zap.Error(err))
					hasErrored = true
					return
				}
			}(v)
		}

		log.Info("Waiting for Units to finish", zap.Int("unit_num", len(nodes)))
		wg.Wait()

		if hasErrored {
			tx.Rollback()
			log.Fatal("One or More Routines failed, aborting migration")
			return
		}

		err := migGraph.MarkNodesRun(nodes...)
		if err != nil {
			tx.Rollback()
			log.Fatal("Could not mark units as executed: ", zap.Error(err))
			return
		}
		// Check if the graph is still shrinkable
		if migGraph.IsStuck() {
			log.Fatal("Migration got stuck on nodes", zap.Strings("nodes", nodes))
			return
		}

		log.Info("Beginning next round")
	}

	log.Info("Committing Migration")
	err = tx.Commit()
	if err != nil {
		log.Fatal("Could not commit migration: %s", zap.Error(err))
		return
	}
	log.Info("Migration finished")
	return
}
