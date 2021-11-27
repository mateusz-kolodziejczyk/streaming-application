package main

import "database/sql"

//import "gopkg.in/guregu/null.v3"

type User struct {
	ID        int    `json:"ID"`
	Username  string `json:"username"`
	StreamKey string `json:"streamKey"`
}

func createUser(db *sql.DB, u User) error {
	_, err := db.Exec("INSERT INTO users (Username, stream_key) VALUES ($1, $2)", u.Username, u.StreamKey)
	return err
}

func getUserByUsername(db *sql.DB, username string) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE Username = $1", username)
	u := User{}
	err := row.Scan(&u.ID, &u.Username, &u.StreamKey)
	return u, err
}

func getUserByStreamKey(db *sql.DB, streamKey string) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE stream_key = $1", streamKey)
	u := User{}
	err := row.Scan(&u.ID, &u.Username, &u.StreamKey)
	return u, err
}
