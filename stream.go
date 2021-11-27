package main

import (
	"database/sql"
	"errors"
	"time"
)

type Stream struct {
	id int
	startTime time.Time
	endTime sql.NullTime
	userID int
}

func createStream(db *sql.DB, userID int) error{
	startTime := time.Now()
	_, err := db.Exec("INSERT INTO streams (start_time, user_id) VALUES ($1, $2)", startTime, userID)
	return err
}

// getOngoingStream Returns an error if the stream has ended and nil if it exists and hasn't.
func getOngoingStream(db *sql.DB, user User) (Stream, error) {
	row := db.QueryRow("SELECT id, end_time, user_id FROM streams WHERE end_time IS NULL and user_id = $1", user.id)
	stream := Stream{}
	err := row.Scan(&stream.id, &stream.endTime, &stream.userID)

	if err != nil {
		return stream, err
	}

	// If the endtime is valid that means that it exists and should not be updated
	if stream.endTime.Valid{
		err = errors.New("stream has already ended")
	}

	return stream, err
}

func endStream(db *sql.DB, user User) error{
	// First check if stream has an end time already
	// If it does return an error
	stream, err := getOngoingStream(db, user)
	if err != nil {
		return err
	}

	endTime := time.Now()
	_, err = db.Exec(`UPDATE streams 
							SET end_time = $1 
							WHERE id = $2`, endTime, stream.id)
	return err
}
