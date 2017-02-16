package mig // import "iris.arke.works/forum/db/mig"

import (
	"database/sql"
	// Import lib/pg for postgres support
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

// PostgresDialect implements a Postgres compatible interface to perform database migration using the mig toolkit
type PostgresDialect struct {
	db *sql.DB
}

// OpenFromPGConn accepts a opened database connections and wraps it into the PostgresDialect
func OpenFromPGConn(db *sql.DB) *PostgresDialect {
	return &PostgresDialect{
		db: db,
	}
}

// CheckAndLoadTables will determine if the database is reachable and create the migration table
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

// MarkExecuted will put a unit into the migration table unless it's marked as "always_exec: true" or a target unit
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

// GetExecutedUnits returns a list of executed units present in the migration table
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
