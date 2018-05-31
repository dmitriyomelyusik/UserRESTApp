package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

//Postgres is a PostgreSQL database
type Postgres struct {
	db *sql.DB
}

//User is type that contains user data
type User struct {
	ID    string
	Email string
	Name  string
	Info  interface{}
}

func (u User) String() string {
	return fmt.Sprintf("ID%v %v %v \t %v", u.ID, u.Name, u.Email, u.Info)
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

//Users returns all users in DB
func (p Postgres) Users() ([]User, error) {
	var users []User
	rows, err := p.db.Query("SELECT * from users")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var u User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, err
		}
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, nil
}

//userHandler handles /user
func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p, err := NewDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		var u []User
		u, err = p.Users()
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

func newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", userHandler)
	return mux
}

func main() {
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      newRouter(),
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
