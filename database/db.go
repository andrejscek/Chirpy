package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	path string
	mux  *sync.Mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string, debug bool) (*DB, error) {
	db := DB{
		path: filepath.Join(path, "database.json"),
		mux:  &sync.Mutex{},
	}

	if debug {
		os.Remove(db.path)
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbs := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}

		data, err := json.Marshal(dbs)
		if err != nil {
			return err
		}

		err = os.WriteFile(db.path, data, os.ModeExclusive)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbs := DBStructure{}
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return DBStructure{}, err
	}

	return dbs, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, os.ModeExclusive)
	if err != nil {
		return err
	}

	return nil
}
