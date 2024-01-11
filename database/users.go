package database

import "sort"

func (db *DB) CreateUser(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	var id int
	var ids []int
	for k := range dbs.Users {
		ids = append(ids, dbs.Users[k].ID)
	}

	if len(ids) == 0 {
		id = 1
	} else {
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		id = ids[len(ids)-1] + 1
	}

	user := User{
		ID:    id,
		Email: email,
	}

	dbs.Users[id] = user

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
