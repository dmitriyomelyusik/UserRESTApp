package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"./postgres"

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

func userHandler(p *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := p.Users()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		for _, v := range u {
			_, err = fmt.Fprintln(w, v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func userIDHandler(p *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "HI")
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := p.GetUserByID(id)
		if err != nil {
			if fmt.Sprintf("%v", err) == "sql: no rows in result set" {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		_, err = fmt.Fprintln(w, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	vars := getEnvVars()
	conf := fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=%v", vars[PGHOST], vars[PGUSER], vars[DBNAME], vars[PGPASS], vars[SSLMODE])
	p, err := postgres.NewDB(conf)
	if err != nil {
		panic(err)
	}
	err = p.DB.Ping()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler(p)).Methods("GET")
	r.HandleFunc("/user/{id}", userIDHandler(p)).Methods("GET")
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
	for _, value := range vars {
		if value == "" {
			panic("You didn't set all requirement environment variables!")
		}
	}
	return vars
}
