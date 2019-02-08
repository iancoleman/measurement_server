package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

var db *sql.DB
var ignoreInitDbConnError = initDbConn()
var ignoreCreateTablesError = createTables()

const dbFilename = "measurements.db?_busy_timeout=5000"

var dbMutex = &sync.Mutex{}

func WriteDb(query string, params ...interface{}) (sql.Result, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return db.Exec(query, params...)
}

func initDbConn() error {
	var err error
	db, err = sql.Open("sqlite3", dbFilename)
	if err != nil {
		panic(err)
	}
	return err
}

func createTables() error {
	var err error
	if err != nil {
		panic(err)
	}
	// measurement table
	_, err = WriteDb(`CREATE TABLE IF NOT EXISTS measurement (
		id integer primary key,
		measured_unix_time real,
		received_unix_time real,
		ip text
	)`)
	if err != nil {
		panic(err)
	}
	// property table
	_, err = WriteDb(`CREATE TABLE IF NOT EXISTS measurement_property (
		measurement_id integer,
		key text,
		value text
	)`)
	if err != nil {
		panic(err)
	}
	// property indexes
	_, err = WriteDb(`CREATE INDEX IF NOT EXISTS measurement_id_index ON measurement_property (measurement_id)`)
	if err != nil {
		panic(err)
	}
	return nil
}
