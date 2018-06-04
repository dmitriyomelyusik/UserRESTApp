package handlers

import (
	"fmt"
	"net/http"

	"github.com/UserRESTApp/entity"

	"github.com/gorilla/mux"
)

//Database is an interface that provide you access to operations with database
type Database interface {
	Users() ([]entity.User, error)
	GetUserByID(id string) (entity.User, error)
}

//UserHandler handles GET/user method
func UserHandler(p Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := p.Users()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(fmt.Sprint(err)))
			if err != nil {
				fmt.Println(err)
			}
		}
		for _, v := range u {
			_, err = fmt.Fprintln(w, v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(fmt.Sprint(err)))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

//UserIDHandler handles GET/user/{id} method
func UserIDHandler(p Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := p.GetUserByID(id)
		if err != nil {
			if fmt.Sprintf("%v", err) == "User is not found" {
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(fmt.Sprintf("Cannot found user with id %v", id)))
				if err != nil {
					fmt.Println(err)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(fmt.Sprint(err)))
				if err != nil {
					fmt.Println(err)
				}
			}
			return
		}
		_, err = fmt.Fprintln(w, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(fmt.Sprint(err)))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}
