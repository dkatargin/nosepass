package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nosepass/common"
	"os"
	"time"
)

type secret struct {
	id         int
	path       string
	created_at int64
	password   string
}

func _getDb() (*sql.DB, error) {
	// Main func for database checks
	configuration, err := common.Config()
	if err != nil {
		return nil, err
	}

	dbPath := fmt.Sprintf("%s/%s", configuration.DbDir, "store.db")
	tableExist := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.MkdirAll(configuration.DbDir, 0700)
		tableExist = false
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if !tableExist {
		log.Println("Initialize table...")
		_, err = db.Exec("create table `secrets` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `path` VARCHAR(255) NOT NULL UNIQUE, `created_at` INTEGER, `password` TEXT NOT NULL)")
		if err != nil {
			return nil, err
		}
	}
	return db, err
}

func storeKey(path string, ciphertext string) error {
	// Insert encrypted key to database
	db, err := _getDb()
	if err != nil {
		return err
	}

	_, err = db.Exec("insert into secrets (path, created_at, password) values ($1, $2, $3)",
		path, time.Now().Unix(), ciphertext)
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}

func getKey(path string) (string, error) {
	// Retrieve encrypted key from database
	db, err := _getDb()
	if err != nil {
		return "", err
	}
	row := db.QueryRow("select * from secrets where path = $1", path)
	if err != nil {
		return "", err
	}
	s := secret{}
	err = row.Scan(&s.id, &s.path, &s.created_at, &s.password)

	defer db.Close()
	return s.password, nil

}

func listKeys() ([]secret, error) {
	// Query list of key paths
	db, err := _getDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from secrets")
	if err != nil {
		return nil, err
	}

	var secrets []secret
	for rows.Next() {
		s := secret{}
		err = rows.Scan(&s.id, &s.path, &s.created_at, &s.password)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, s)
	}

	defer rows.Close()
	defer db.Close()
	return secrets, err
}

func deleteKey(dstPath string) error {
	// Delete encrypted key from database
	db, err := _getDb()
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from secrets where path = $1", dstPath)
	defer db.Close()
	return err
}
