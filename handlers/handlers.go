package handlers

import (
	"fmt"
	"net/http"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/errors"

	"github.com/gorilla/mux"
)

//Response message when user wasn't found
const (
	NotFound = "Cannot found user with id "
)

//Server uses controller methods to work with them together with http methods
type Server struct {
	Controller controller.User
}

//UserHandler handles GET/user method
func (s Server) UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.Controller.GetUsers()
		if err != nil {
			occuredError(err, w, "")
			return
		}
		for _, v := range u {
			_, err = fmt.Fprintln(w, v)
			if err != nil {
				occuredError(err, w, "")
				return
			}
		}
	}
}

//UserIDHandler handles GET/user/{id} method
func (s Server) UserIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := s.Controller.GetUserByID(id)
		if err != nil {
			occuredError(err, w, id)
			return
		}
		_, err = fmt.Fprintln(w, u)
		if err != nil {
			occuredError(err, w, id)
			return
		}
	}
}

//NewRouter returns router with configurated and handled pathes
func NewRouter(s Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user", s.UserHandler()).Methods("GET")
	r.HandleFunc("/user/", s.UserHandler()).Methods("GET")
	r.HandleFunc("/user/{id}", s.UserIDHandler()).Methods("GET")
	return r
}

func occuredError(err error, w http.ResponseWriter, id string) {
	switch err {
	case errors.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte(fmt.Sprintf("%v %v", NotFound, id)))
		if err != nil {
			fmt.Println(err)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprint(err)))
		if err != nil {
			fmt.Println(err)
		}

	}
}
