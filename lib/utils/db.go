package utils

import (
	"database/sql"
	"fmt"
	"keypass/lib/stores"
	"time"
)

func SetPragmaWAL() error {
	dbconn := GetLiteConn()
	defer dbconn.Close()
	_, err := dbconn.Exec("PRAGMA journal_mode=WAL")
	return err
}

func GetLiteConn() *sql.DB {
	db, err := sql.Open("sqlite3", stores.Config.DbPath)
	if err != nil {
		Logging("GetLiteConn1", err.Error())
		panic(err)
	}
	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(50)
	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetConnMaxLifetime(time.Hour)
	return db
}

// QueryRow function get single row
func QueryRow(query string, args ...interface{}) *sql.Row {
	Logging("QueryRow", fmt.Sprintf("%s %v", query, args))
	dbconn := GetLiteConn()
	defer dbconn.Close()
	row := dbconn.QueryRow(query, args...)
	return row
}

// QueryDb function to query sqlite
func QueryDb(query string, args ...interface{}) (*sql.Rows, error) {
	Logging("QueryDb", fmt.Sprintf("%s %v", query, args))
	dbconn := GetLiteConn()
	defer dbconn.Close()
	rows, err := dbconn.Query(query, args...)
	return rows, err
}

func TrxDb(query string, args ...interface{}) (sql.Result, error) {
	dbconn := GetLiteConn()
	defer dbconn.Close()
	var sqlRes sql.Result
	trashSQL, err := dbconn.Prepare(query)
	if err != nil {
		return sqlRes, err
	}
	tx, err := dbconn.Begin()
	if err != nil {
		return sqlRes, err
	}
	sqlRes, err = tx.Stmt(trashSQL).Exec(args...)
	if err != nil {
		tx.Rollback()
		return sqlRes, err
	}
	tx.Commit()
	return sqlRes, err
}

// InsertDb function to insert sqlite
func InsertDb(query string, args ...interface{}) (int64, error) {
	Logging("InsertDb", fmt.Sprintf("%s %v", query, args))
	res, err := TrxDb(query, args...)
	if err != nil {
		Logging("InsertDb1", fmt.Sprintf("%q: %s", err, query))
		return -1, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		Logging("InsertDb2", fmt.Sprintf("%q: %s", err, query))
		return -1, err
	}
	return lastID, nil
}

// ExecDb to update or delete sqlite
func ExecDb(query string, args ...interface{}) (int64, error) {
	Logging("ExecDb", fmt.Sprintf("%s %v", query, args))
	res, err := TrxDb(query, args...)
	if err != nil {
		Logging("ExecDb1", fmt.Sprintf("%q: %s", err, query))
		return -1, err
	}
	if err != nil {
		Logging("ExecDb2", fmt.Sprintf("%q: %s", err, query))
		return -1, err
	}
	affectedRows, err := res.RowsAffected()
	return affectedRows, err
}
