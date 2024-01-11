package database

import "sort"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (db *DB) CreateUser(email string, pwd []byte) (User, error) {
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
		ID:       id,
		Email:    email,
		Password: pwd,
	}

	dbs.Users[id] = user

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, u := range dbs.Users {
		if u.Email == email {
			return u, nil
		}
	}

	return User{}, err
}
