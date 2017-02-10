package mig // import "iris.arke.works/forum/db/mig"

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const pgMigTable = `CREATE TABLE IF NOT EXISTS vape_migrations (
	name    varchar(1024)   NOT NULL,
	type    varchar(1024)   NOT NULL,
	executed_on timestamptz NOT NULL    DEFAULT (now() AT TIME ZONE 'utc'),

	PRIMARY KEY(name)
);

DROP INDEX IF EXISTS vape_type_index;
DROP INDEX IF EXISTS vape_name_type_index;
CREATE INDEX vape_type_index ON vape_migrations(type);
CREATE INDEX vape_name_type_index ON vape_migrations(name, type);

--`

const markExecutedQuery = `INSERT INTO vape_migrations (name, type) VALUES
($1, $2)
;`

const getExecutedQuery = `SELECT name FROM vape_migrations WHERE type=$1;`

type PostgresDialect struct {
	db *sql.DB
}

func OpenFromPGConn(db *sql.DB) *PostgresDialect {
	return &PostgresDialect{
		db: db,
	}
}

// CheckTablesAndTables will determine if the migration table is present in the
// database and load
func (d *PostgresDialect) CheckAndLoadTables() error {
	err := d.db.Ping()
	if err != nil {
		return err
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(pgMigTable)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *PostgresDialect) MarkExecuted(unit Unit) error {
	if unit.Type == UnitTypeVirtualTarget {
		return nil
	}
	if unit.AlwaysExec {
		return nil
	}
	_, err := d.db.Exec(markExecutedQuery, unit.Name, string(unit.Type))
	return err
}

func (d *PostgresDialect) GetExecutedUnits() ([]string, error) {
	rows, err := d.db.Query(getExecutedQuery, string(UnitTypeMigration))
	if err != nil {
		return nil, err
	}
	var units = []string{}
	for rows.Next() {
		var unit string
		err = rows.Scan(&unit)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}
