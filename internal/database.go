package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var db *sql.DB

func initDB() {
	if workDir, err := os.Getwd(); err == nil {
		if db, err = sql.Open("sqlite3", workDir+"\\database\\instagram.db"); err == nil {
			if err := createTables(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func createTables() error {
	createTables := `CREATE TABLE IF NOT EXISTS targetUsers (
    				 id 			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
					 instagram_id	BIGINT NOT NULL,
					 instagram_name NVARCHAR(255) NOT NULL,
					 deleted		BOOLEAN DEFAULT 0);
					 
					 CREATE TABLE IF NOT EXISTS followers(
					 id			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
					 owner		INTEGER NOT NULL,
					 add_date	date NOT NULL,
					 deleted	BOOLEAN DEFAULT 0,
					 FOREIGN KEY (owner) REFERENCES targetUsers(id));
					 
					 CREATE TABLE IF NOT EXISTS followings(
					 id			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
					 owner		INTEGER NOT NULL,
					 add_date	date NOT NULL,
					 deleted	BOOLEAN DEFAULT 0,
					 FOREIGN KEY (owner) REFERENCES targetUsers(id));`

	if _, err := db.Exec(createTables); err != nil {
		return err
	}
	return nil
}
