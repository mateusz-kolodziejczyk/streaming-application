package main

import (
	"database/sql"
	"errors"
	"gopkg.in/guregu/null.v3"
	"time"
)

type Stream struct {
	ID        int       `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   null.Time `json:"endTime"`
	UserID    int       `json:"userID"`
}

func createStream(db *sql.DB, userID int) error {
	startTime := time.Now()
	_, err := db.Exec("INSERT INTO streams (start_time, user_id) VALUES ($1, $2)", startTime, userID)
	return err
}

// getOngoingStream Returns an error if the stream has ended and nil if it exists and hasn't.
func getOngoingStream(db *sql.DB, u User) (Stream, error) {
	row := db.QueryRow("SELECT ID, end_time, user_id FROM streams WHERE end_time IS NULL and user_id = $1", u.ID)
	stream := Stream{}
	err := row.Scan(&stream.ID, &stream.EndTime, &stream.UserID)

	if err != nil {
		return stream, err
	}

	// If the endtime is valid that means that it exists and should not be updated
	if stream.EndTime.Valid {
		err = errors.New("stream has already ended")
	}

	return stream, err
}

func endStream(db *sql.DB, u User) error {
	// First check if stream has an end time already
	// If it does return an error
	stream, err := getOngoingStream(db, u)
	if err != nil {
		return err
	}

	endTime := time.Now()
	_, err = db.Exec(`UPDATE streams 
							SET end_time = $1 
							WHERE ID = $2`, endTime, stream.ID)
	return err
}
