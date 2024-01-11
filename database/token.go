package database

import "time"

type Token struct {
	ID         string    `json:"id"`
	RevokeTime time.Time `json:"revoke_time"`
}

func (db *DB) CheckRevoked(token string) (Token, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}

	revoked, ok := dbs.Revoked[token]
	if !ok {
		return Token{}, err
	}

	return revoked, nil
}

func (db *DB) RevokeToken(token string) (Token, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}

	to_revoke := Token{
		ID:         token,
		RevokeTime: time.Now(),
	}

	dbs.Revoked[token] = to_revoke

	err = db.writeDB(dbs)
	if err != nil {
		return Token{}, err
	}

	return to_revoke, nil
}
