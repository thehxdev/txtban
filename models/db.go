package models

import (
	"database/sql"
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
    content VARCHAR(4096) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uid INT NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(id)
);`

type Conn struct {
	DB *sql.DB
}

func (c *Conn) MigrateDB() {
	_, err := c.DB.Exec(dbSchema)
	if err != nil {
		panic(err)
	}
}
