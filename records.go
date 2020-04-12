package main

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

func connectDatabase(sqldataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", sqldataSource)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func createTableStatement(tableName string) []string {
	return []string{`
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		time TIMESTAMP NOT NULL,
		value FLOAT8 NOT NULL,
		UNIQUE (time)
	);
	`}
}

func setupTables(db *sql.DB, tableName string, queryBuilder func(string) []string) error {
	setupTrans := func(trans *sql.Tx, table string, makeStmt func(string) []string) error {
		for _, stmt := range makeStmt(table) {
			if _, err := trans.Exec(stmt); err != nil {
				return errors.Wrapf(err, "failed to execute statement %s", stmt)
			}
		}
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := setupTrans(tx, tableName, queryBuilder); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func setupDB(sqldataSource string) (*sql.DB, error) {
	var lastErr error
	// the db might be unavailable for some time, so let try for 5 times for 25 seconds
	for i := 0; i < 5; i++ {

		db, err := connectDatabase(sqldataSource)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		if err = setupTables(db, "raised", createTableStatement); err != nil {
			db.Close()
			return nil, err
		}
		return db, err
	}
	return nil, lastErr
}

func writeValue(db *sql.DB, t time.Time, v float64) error {
	_, err := db.Exec("INSERT INTO raised VALUES($1,$2)", t.UTC(), v)
	return err
}
