package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	//"io"
	"log"
	"net/http"
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

func getStreamURL() {

}

func postUser() {

}

func startStreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["Username"]
	user, err := getUserByUsername(app.DB, username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	_, err = getOngoingStream(app.DB, user)

	// If there is no error then an ongoing stream exists for the user
	if err == nil {
		http.Error(w, "Stream already running", http.StatusConflict)
		return
	}
	// Probe the RTMP stream to check if it's running.
	if !probeRTMPStream(user.StreamKey, WinServerAddress){
		http.Error(w, "Origin stream not found", http.StatusNotFound)
		return
	}
	go startHLSStream(3, fmt.Sprintf("%s\\stream\\%s", app.Path, username), user.StreamKey)

	w.WriteHeader(http.StatusOK)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var u User

	// Turn the request body into a user struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = createUser(app.DB, u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Load environmental variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app.Initialize()

	streamPath := fmt.Sprintf("%s\\stream", app.Path)
	//go startHLSStream(3, streamPath, "cool")

	app.Router.HandleFunc("/api/user", addUserHandler).Methods("POST")
	app.Router.HandleFunc("/api/stream/{Username}", startStreamHandler).Methods("POST")
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
