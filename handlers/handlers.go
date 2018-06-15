// Package handlers has methods to operate with client requests via http methods.
// All of data are encoded via json and wroten to response.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/UserRESTApp/entity"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/errors"

	"github.com/gorilla/mux"
)

// Server uses controller methods to work with them together with http methods
type Server struct {
	Controller controller.User
}

// UsersHandler handles GET /user/ method
func (s Server) UsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.Controller.GetUsers()
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, u)
	}
}

// UserHandler handles GET /user/{id} method
func (s Server) UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := s.Controller.GetUser(id)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, u)
	}
}

// CreateUserHandler handles POST /user/ method
func (s Server) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u entity.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			jsonError(w, err)
			return
		}
		id, err := s.Controller.CreateUser(u)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, id)
	}
}

// PutUserHandler handles PUT /user/{id} method
func (s Server) PutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u entity.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			jsonError(w, err)
			return
		}
		vars := mux.Vars(r)
		u.ID = vars["id"]
		if err := s.Controller.UpdateUser(u); err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// DeleteUserHandler handles DELETE /user/{id} method
func (s Server) DeleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if err := s.Controller.DeleteUser(id); err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// PatchUserHandler handles PATCH /user/{id} method
func (s Server) PatchUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		changes := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
			jsonError(w, err)
			return
		}
		if err := s.Controller.UpdateUserFields(changes, id); err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// NewRouter returns router with configurated and handled pathes
func NewRouter(s Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/", s.UsersHandler()).Methods("GET")
	r.HandleFunc("/user/{id}", s.UserHandler()).Methods("GET")
	r.HandleFunc("/user/", s.CreateUserHandler()).Methods("POST")
	r.HandleFunc("/user/{id}", s.PutUserHandler()).Methods("PUT")
	r.HandleFunc("/user/{id}", s.DeleteUserHandler()).Methods("DELETE")
	r.HandleFunc("/user/{id}", s.PatchUserHandler()).Methods("PATCH")
	return r
}

func jsonError(w http.ResponseWriter, err error) {
	myErr, ok := err.(errors.Error)
	if !ok {
		myErr = errors.Error{
			Code:    "UnknownError",
			Message: err.Error(),
		}
	}
	switch myErr.Code {
	case errors.UserNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errors.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case errors.UserExists:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResponse(w, myErr)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		log.Println(err)
	}
}
