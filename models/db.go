package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
)

const dbSchema = `
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    phash VARCHAR(60) NOT NULL UNIQUE,
    authKey VARCHAR(92) NOT NULL UNIQUE
);

CREATE TABLE txts (
    id VARCHAR(16) NOT NULL UNIQUE,
    name VARCHAR(32) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uid INT NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(id)
);`

type DB struct {
	Write *sql.DB
	Read  *sql.DB
}

func (d *DB) MigrateDB() {
	_, err := d.Write.Exec(dbSchema)
	if err != nil {
		panic(err)
	}
}

func (d *DB) SetupSqliteDB(path string) error {
	connUrlParams := make(url.Values)
	connUrlParams.Add("_txlock", "immediate")
	connUrlParams.Add("_journal_mode", "WAL")
	connUrlParams.Add("_busy_timeout", "5000")
	connUrlParams.Add("_synchronous", "NORMAL")
	// connUrlParams.Add("_cache_size", "1000000000")
	connUrlParams.Add("_foreign_keys", "true")
	connUrl := fmt.Sprintf("file:%s?%s", path, connUrlParams.Encode())

	writeDB, err := sql.Open("sqlite3", connUrl)
	if err != nil {
		return err
	}
	writeDB.SetMaxOpenConns(1)
	d.Write = writeDB

	readDB, err := sql.Open("sqlite3", connUrl)
	if err != nil {
		return err
	}
	d.Read = readDB
	readDB.SetMaxOpenConns(max(1, runtime.NumCPU()))

	return nil
}
