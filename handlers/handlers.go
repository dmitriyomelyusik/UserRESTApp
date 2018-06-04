package handlers

import (
	"fmt"
	"net/http"

	"../entity"

	"github.com/gorilla/mux"
)

//UserHandler handles GET/user method
func UserHandler(p entity.Database) http.HandlerFunc {
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

//UserIDHandler handles GET/user/{id} method
func UserIDHandler(p entity.Database) http.HandlerFunc {
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
