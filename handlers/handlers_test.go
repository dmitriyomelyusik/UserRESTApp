// Package handlers_test has a set of tests needed to check correctness of wroten handlers.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
	"github.com/UserRESTApp/postgres"
	"github.com/gavv/httpexpect"

	_ "github.com/lib/pq"
)

var (
	ts *httptest.Server
)

func TestMain(m *testing.M) {
	p, err := postgres.NewDB("user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	r := NewRouter(Server{Controller: controller.User{DB: p}})
	ts = httptest.NewServer(r)
	defer ts.Close()
	code := m.Run()
	os.Exit(code)
}

func TestHandlers_GetHandler(t *testing.T) {
	users := []entity.User{
		{Name: "gethandler_test1", Email: "gethandler_test1", Info: "gethandler_test1"},
		{Name: "gethandler_test2", Email: "gethandler_test2", Info: "gethandler_test2"},
		{Name: "gethandler_test3", Email: "gethandler_test3", Info: "gethandler_test3"},
	}
	e := httpexpect.New(t, ts.URL)
	var userIDs []string
	for i := range users {
		rawID := e.POST("/user/").WithJSON(users[i]).Expect().Status(http.StatusOK).Body().Raw()
		var id string
		err := json.Unmarshal([]byte(rawID), &id)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
		users[i].ID = id
	}
	for i := range users {
		u, err := json.Marshal(users[i])
		if err != nil {
			t.Fatal(err)
		}
		e.GET("/user/").Expect().Status(http.StatusOK).Body().Contains(string(u))
	}
	e.GET("/user").Expect().Status(http.StatusNotFound)
	for i := range userIDs {
		e.DELETE("/user/" + userIDs[i]).Expect().Status(http.StatusOK)
	}
}

func TestHandlers_GetIDHandler(t *testing.T) {
	users := []entity.User{
		{Name: "getidhandler_test1", Email: "getidhandler_test1", Info: "getidhandler_test1"},
		{Name: "getidhandler_test2", Email: "getidhandler_test2", Info: "getidhandler_test2"},
		{Name: "getidhandler_test3", Email: "getidhandler_test3", Info: "getidhandler_test3"},
	}
	e := httpexpect.New(t, ts.URL)
	var userIDs []string
	for i := range users {
		rawID := e.POST("/user/").WithJSON(users[i]).Expect().Status(http.StatusOK).Body().Raw()
		var id string
		err := json.Unmarshal([]byte(rawID), &id)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
		users[i].ID = id
	}
	for i := range users {
		u, err := json.Marshal(users[i])
		if err != nil {
			t.Fatal(err)
		}
		e.GET("/user/" + userIDs[i]).Expect().Status(http.StatusOK).Body().Equal(string(u) + "\n")
	}
	e.GET("/user/getidhandler_wrongID").Expect().Status(http.StatusNotFound)
	for i := range userIDs {
		e.DELETE("/user/" + userIDs[i]).Expect().Status(http.StatusOK)
	}
}

func TestHandlers_DeleteHandler(t *testing.T) {
	users := []entity.User{
		{Name: "deletehandler_test1", Email: "deletehandler_test1", Info: "deletehandler_test1"},
		{Name: "deletehandler_test2", Email: "deletehandler_test2", Info: "deletehandler_test2"},
		{Name: "deletehandler_test3", Email: "deletehandler_test3", Info: "deletehandler_test3"},
	}
	e := httpexpect.New(t, ts.URL)
	var userIDs []string
	for i := range users {
		rawID := e.POST("/user/").WithJSON(users[i]).Expect().Status(http.StatusOK).Body().Raw()
		var id string
		err := json.Unmarshal([]byte(rawID), &id)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
		users[i].ID = id
	}
	for i := range userIDs {
		e.DELETE("/user/" + userIDs[i]).Expect().Status(http.StatusOK)
	}
	wrongIDs := []string{
		"deletehandler_wrongID1",
		"deletehandler_wrongID2",
		"deletehandler_wrongID3",
		"deletehandler_wrongID4",
		"deletehandler_wrongID5",
	}
	for i := range wrongIDs {
		myErr := errors.Error{Code: errors.UserNotFound, Message: "cannot find user in database with id " + wrongIDs[i]}
		b, err := json.Marshal(myErr)
		if err != nil {
			t.Fatal(err)
		}
		e.DELETE("/user/" + wrongIDs[i]).Expect().Status(http.StatusNotFound).Body().Equal(string(b) + "\n")
	}
}

func TestHandlers_CreateHandler(t *testing.T) {
	users := []entity.User{
		{Name: "createhandler_test1", Email: "createhandler_test1", Info: "createhandler_test1"},
		{Name: "createhandler_test2", Email: "createhandler_test2", Info: "createhandler_test2"},
		{Name: "createhandler_test3", Email: "createhandler_test3", Info: "createhandler_test3"},
	}
	e := httpexpect.New(t, ts.URL)
	var userIDs []string
	for i := range users {
		rawID := e.POST("/user/").WithJSON(users[i]).Expect().Status(http.StatusOK).Body().Raw()
		var id string
		err := json.Unmarshal([]byte(rawID), &id)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
		users[i].ID = id
	}
	//TODO: add some tests
	for i := range userIDs {
		e.DELETE("/user/" + userIDs[i]).Expect().Status(http.StatusOK)
	}
}
