package main

import (
	"database/sql"
	"time"
)

type Stream struct {
	id int
	startTime time.Time
	endTime sql.NullTime
	userID int
}

// createStream Add a stream linked to userID with start time set to the current time to the database.
func createStream(db *sql.DB, userID int) error{
	startTime := time.Now()
	_, err := db.Exec("INSERT INTO streams (start_time, user_id) VALUES ($1, $2)", startTime, userID)
	return err
}