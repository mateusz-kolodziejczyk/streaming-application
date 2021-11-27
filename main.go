package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	//"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const StreamServerDirectory string = "C:\\Users\\MK\\GolandProjects\\streamingApplication\\build"
const ServerAddress string = "127.0.0.1"
const LocalServerPath string = "build"
const WinServerAddress string = "192.168.0.66"
const StreamDirectory string = "stream"
const ServerDirectory string = "/usr/local/nginx/html"
const TimeoutMicroSeconds int = 5000000

var outputResolutions [][]int = [][]int{{1920, 1080}, {1280, 720}, {854, 480}, {640, 360}}

// MaxBitrate Bitrate in megabits
const MaxBitrate float64 = 5


var app App

func getStreamURL(){

}

func postUser(){

}

func endStream(db *sql.DB, streamID int) error{
	// First check if stream has an end time already
	// If it does return an error
	err := streamOngoing(db, streamID)
	if err != nil {
		return err
	}

	endTime := time.Now()
	_, err = db.Exec(`UPDATE streams 
							SET end_time = $1 
							WHERE id = $2`, endTime, streamID)
	return err
}

func startStreamHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	username := vars["username"]
	_, err := getUser(app.DB, username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

}

// streamOngoing Returns an error if the stream has ended and nil if it exists and hasn't.
func streamOngoing(db *sql.DB, streamID int) error {
	row := db.QueryRow("SELECT end_time FROM streams WHERE id = $1", streamID)
	stream := Stream{}
	err := row.Scan(&stream.endTime)

	if err != nil {
		return err
	}

	// If the endtime is valid that means that it exists and should not be updated
	if stream.endTime.Valid{
		err = errors.New("stream has already ended")
	}

	return err
}


func main() {
	// Load environmental variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app.Initialize()
	// Get the path to server
	localpath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	streamPath := fmt.Sprintf("%s\\%s\\stream", localpath, LocalServerPath)
	//go startHLSStream(3, streamPath, "cool")

	app.Router.HandleFunc("/api/stream/{username}", startStreamHandler).Methods("GET")
	app.Router.PathPrefix("/stream/").Handler(http.StripPrefix("/stream/", http.FileServer(http.Dir(filepath.Join(streamPath)))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},   // All origins
		AllowedMethods: []string{"GET"}, // Allowing only get, just an example
	})
	srv := &http.Server{
		Handler: c.Handler(app.Router),
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
