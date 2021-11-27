package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"os"
)

// Code from https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

type App struct {
	Router *mux.Router
	DB   *sql.DB
	Path string
}

func (a *App) Initialize() {
	var err error
	a.DB, err = connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := a.DB.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")
	// Get the path to server
	a.Path, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	a.Path = fmt.Sprintf("%s\\%s", app.Path, os.Getenv("SERVER_DIRECTORY"))
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) { }