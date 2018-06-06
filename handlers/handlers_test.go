package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
	. "github.com/UserRESTApp/handlers"
	"github.com/UserRESTApp/postgres"

	_ "github.com/lib/pq"
)

var (
	p  *postgres.Postgres
	r  *mux.Router
	ts *httptest.Server
)

func TestMain(m *testing.M) {
	var err error
	p, err = postgres.NewDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		fmt.Printf("Cannot open database: %v", err)
	}
	r = NewRouter(Server{Controller: controller.User{DB: p}})
	code := m.Run()
	os.Exit(code)
}

func TestUserHanlder(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user")
	if err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(res.Body)
	var user entity.User
	for decoder.More() {
		err := decoder.Decode(&user)
		if err != nil {
			t.Fatal("Bad user perfomance")
		}
	}
}

func TestErrorsMessages(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	id := []string{"@-1", "@-2", "@-3", "@@@"}

	for _, v := range id {
		res, err := http.Get(ts.URL + "/user/" + v)
		if err != nil {
			t.Fatal(err)
		}

		decoder := json.NewDecoder(res.Body)
		var code errors.ErrCode
		err = decoder.Decode(&code)
		fmt.Println(v)
		if code != errors.UserNotFound {
			t.Fatalf("Wrong err code: expected %v, got %v", errors.UserNotFound, code)
		}
	}
}
