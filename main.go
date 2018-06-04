package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/UserRESTApp/handlers"
	"github.com/UserRESTApp/postgres"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Environment variables that needs to open database
const (
	PGHOST  = "PGHOST"
	DBNAME  = "DBNAME"
	PGUSER  = "PGUSER"
	PGPASS  = "PGPASS"
	SSLMODE = "SSLMODE"
)

func main() {
	vars := getEnvVars()
	conf := fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=%v", vars[PGHOST], vars[PGUSER], vars[DBNAME], vars[PGPASS], vars[SSLMODE])
	p, err := postgres.NewDB(conf)
	if err != nil {
		fmt.Println("Fatal to open database. Check environment variables.")
		panic(err)
	}
	err = p.DB.Ping()
	if err != nil {
		fmt.Println("Fatal to connet to database. Check environment variables.")
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/user", handlers.UserHandler(p)).Methods("GET")
	r.HandleFunc("/user/{id}", handlers.UserIDHandler(p)).Methods("GET")
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      r,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnvVars() map[string]string {
	vars := make(map[string]string)
	vars[PGUSER] = os.Getenv(PGUSER)
	vars[PGPASS] = os.Getenv(PGPASS)
	vars[PGHOST] = os.Getenv(PGHOST)
	vars[DBNAME] = os.Getenv(DBNAME)
	vars[SSLMODE] = os.Getenv(SSLMODE)
	isOK := true
	for key, value := range vars {
		if value == "" {
			fmt.Printf("You didn't set environment variable: %v\n", key)
			isOK = false
		}
	}
	if !isOK {
		panic("Set all environment variables!")
	}
	return vars
}
