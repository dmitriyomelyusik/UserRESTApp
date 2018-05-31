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

type user struct {
	id    string
	email string
	name  string
	info  interface{}
}

func (u user) String() string {
	return fmt.Sprintf("ID%v %v %v \t %v", u.id, u.name, u.email,
		u.info)
}

var db *sql.DB

func initDB(conf string) {
	var err error
	db, err = sql.Open("postgres", conf)
	if err != nil {
		log.Fatalf("Fatal to open db: %v", err)
	}
}

func getUsersFromDB() []user {
	var users []user
	rows, err := db.Query("SELECT * from users")
	if err != nil {
		log.Fatalf("Fatal to read users from db: %v", err)
	}
	for rows.Next() {
		var id, email, name string
		var info interface{}
		var infoData []byte
		err = rows.Scan(&id, &email, &name, &infoData)
		if err != nil {
			log.Fatalf("Fatal to scan user: %v", err)
		}
		err = json.Unmarshal(infoData, &info)
		if err != nil {
			log.Fatalf("Fatal to unmarshal info: %v", err)
		}
		users = append(users, user{id: id, email: email, name: name, info: info})
	}
	return users
}

//HandleGet handles get method
func HandleGet(w http.ResponseWriter, r *http.Request) {
	users := getUsersFromDB()
	var err error
	for _, v := range users {
		_, err = fmt.Fprintln(w, v)
		if err != nil {
			log.Fatalf("Fatal to show users: %v", err)
		}
	}
}

func newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/GET/user", HandleGet)
	return mux
}

func main() {
	initDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
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
