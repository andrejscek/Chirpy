package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.Mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	var id int
	var ids []int
	for k := range dbs.Chirps {
		ids = append(ids, dbs.Chirps[k].ID)
	}

	if len(ids) == 0 {
		id = 1
	} else {
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		id = ids[len(ids)-1] + 1
	}

	chirp := Chirp{
		ID:   id,
		Body: body,
	}

	dbs.Chirps[id] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	for _, c := range dbs.Chirps {
		if c.ID == id {
			return c, nil
		}
	}

	return Chirp{}, err
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var keys []int
	for k := range dbs.Chirps {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return []Chirp{}, nil
	} else {

		sortFunc := func(i, j int) bool { return dbs.Chirps[keys[i]].ID < dbs.Chirps[keys[j]].ID }
		sort.Slice(keys, sortFunc)
		chirps := make([]Chirp, len(keys))
		for i := 0; i < len(keys); i++ {
			chirps[i] = dbs.Chirps[keys[i]]
		}

		return chirps, nil
	}
}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbs := DBStructure{
			Chirps: make(map[int]Chirp),
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
