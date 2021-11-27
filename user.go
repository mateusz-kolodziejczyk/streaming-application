package main

import "database/sql"

type User struct {
	id int
	username sql.NullString
	streamKey string
}

func createUser(db *sql.DB, username string, streamKey string) error{
	_, err := db.Exec("INSERT INTO users (username, stream_key) VALUES ($1, $2)", username, streamKey)
	return err
}

func getUserByUsername(db *sql.DB, username string) (User, error){
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	user := User{}
	err := row.Scan(&user.id, &user.username, &user.streamKey)
	return user, err
}

func getUserByStreamKey(db *sql.DB, streamKey string) (User, error){
	row := db.QueryRow("SELECT * FROM users WHERE stream_key = $1", streamKey)
	user := User{}
	err := row.Scan(&user.id, &user.username, &user.streamKey)
	return user, err
}