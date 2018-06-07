package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/UserRESTApp/entity"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/errors"

	"github.com/gorilla/mux"
)

//Server uses controller methods to work with them together with http methods
type Server struct {
	Controller controller.User
}

//UserHandler handles GET /user/ method
func (s Server) UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.Controller.GetUsers()
		if err != nil {
			jsonError(w, err)
			return
		}
		for _, v := range u {
			jsonResponse(w, v)
		}
	}
}

//UserIDHandler handles GET /user/{id} method
func (s Server) UserIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		u, err := s.Controller.GetUserByID(id)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, u)
	}
}

//PostUserHandler handles POST /user/ method
func (s Server) PostUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u entity.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			jsonResponse(w, err)
		}
		err = s.Controller.PostUser(u)
		if err != nil {
			jsonError(w, err)
		}
	}
}

//PutUserByIDHandler handles PUT /user/{id} method
func (s Server) PutUserByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u entity.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			jsonResponse(w, err)
		}
		vars := mux.Vars(r)
		id := vars["id"]
		err = s.Controller.PutUserByID(u, id)
		if err != nil {
			jsonError(w, err)
		}
	}
}

//DeleteUserByIDHandler handles DELETE /user/{id} method
func (s Server) DeleteUserByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := s.Controller.DeleteUserByID(id)
		if err != nil {
			jsonError(w, err)
		}
	}
}

//NewRouter returns router with configurated and handled pathes
func NewRouter(s Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/", s.UserHandler()).Methods("GET")
	r.HandleFunc("/user/{id}", s.UserIDHandler()).Methods("GET")
	r.HandleFunc("/user/", s.PostUserHandler()).Methods("POST")
	r.HandleFunc("/user/{id}", s.PutUserByIDHandler()).Methods("PUT")
	r.HandleFunc("/user/{id}", s.DeleteUserByIDHandler()).Methods("DELETE")
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
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(data)
}
