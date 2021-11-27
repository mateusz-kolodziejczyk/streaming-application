package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
)

// Code from https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

type App struct {
	Router *mux.Router
	DB *sql.DB
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
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) { }