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
