package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(db.path)

			if err != nil {
				return err
			}

			dbData := DBStructure{
				Chirps: map[int]Chirp{},
				Users:  map[int]User{},
			}

			err = db.writeDB(dbData)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	dbDataBytes, err := os.ReadFile(db.path)

	if err != nil {
		return DBStructure{}, err
	}

	dbData := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}

	err = json.Unmarshal(dbDataBytes, &dbData)

	if err != nil {
		return DBStructure{}, err
	}

	return dbData, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	dbDataBytes, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dbDataBytes, 0644)

	return err
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) *DB {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	if err != nil {
		panic(err)
	}

	return db
}
