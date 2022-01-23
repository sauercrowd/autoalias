package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	historyDb *sql.DB
	aliasDb   *sql.DB
}

func getSqliteFilepath(name string) string {
	storePath, ok := os.LookupEnv("HOME")
	if !ok {
		log.Fatal("Cannot determine home directory, please set the HOME envionment variable")
	}
	return path.Join(storePath, fmt.Sprintf(".autoalias.%s.sqlite", name))
}

func New() (*Store, error) {
	historyDb, err := sql.Open("sqlite3", "file:"+getSqliteFilepath("history"))
	if err != nil {
		return nil, err
	}
	aliasDb, err := sql.Open("sqlite3", "file:"+getSqliteFilepath("alias"))
	if err != nil {
		return nil, err
	}
	if _, err := historyDb.Exec(`CREATE TABLE IF NOT EXISTS history
		(id INTEGER PRIMARY KEY AUTOINCREMENT, command_array TEXT, used TEXT)
	`); err != nil {
		return nil, err
	}

	if _, err := aliasDb.Exec(`CREATE TABLE IF NOT EXISTS aliases
		(id INTEGER PRIMARY KEY AUTOINCREMENT, alias TEXT, command TEXT, created TEXT, last_used TEXT)
	`); err != nil {
		return nil, err
	}
	return &Store{
		historyDb: historyDb,
		aliasDb:   aliasDb,
	}, nil
}

func JsonArrayToArray(jsonArr string) ([]string, error) {
	var arr []string
	if err := json.Unmarshal([]byte(jsonArr), &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func ArrayToJsonArray(arr []string) (string, error) {
	content, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(content[:]), nil
}
