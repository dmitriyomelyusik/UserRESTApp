package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./entity"
	"./postgres"

	_ "github.com/lib/pq"
)

//userHandler handles /user
func userHandler(p *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var u []entity.User
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
}

func newRouter(p *postgres.Postgres) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", userHandler(p))
	return mux
}

func main() {
	p, err := postgres.NewDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		panic(err)
	}
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      newRouter(p),
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
