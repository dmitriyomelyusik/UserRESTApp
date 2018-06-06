package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/errors"

	"github.com/gorilla/mux"
)

//Server uses controller methods to work with them together with http methods
type Server struct {
	Controller controller.User
}

//UserHandler handles GET/user method
func (s Server) UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		u, err := s.Controller.GetUsers()
		if err != (errors.Error{}) {
			jsonError(encoder, err)
			return
		}
		for _, v := range u {
			encoder.Encode(v)
		}
	}
}

//UserIDHandler handles GET/user/{id} method
func (s Server) UserIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := s.Controller.GetUserByID(id)
		if err != (errors.Error{}) {
			jsonError(encoder, errors.Error(err))
			return
		}
		encoder.Encode(u)
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

func jsonError(encoder *json.Encoder, err errors.Error) {
	encoder.Encode(err.Code)
	encoder.Encode(err.Message)
	if err.Info != nil {
		encoder.Encode(err.Info)
	}
}
