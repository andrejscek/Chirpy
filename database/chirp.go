package database

import "sort"

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(body string, author int) (Chirp, error) {
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
		ID:       id,
		AuthorID: author,
		Body:     body,
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
