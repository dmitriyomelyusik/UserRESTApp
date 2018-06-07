package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/UserRESTApp/controller"
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
	. "github.com/UserRESTApp/handlers"
	"github.com/UserRESTApp/postgres"

	_ "github.com/lib/pq"
)

var (
	ts *httptest.Server
)

func TestMain(m *testing.M) {
	p, err := postgres.NewDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	r := NewRouter(Server{Controller: controller.User{DB: p}})
	ts = httptest.NewServer(r)
	defer ts.Close()
	code := m.Run()
	os.Exit(code)
}

func TestUserHanlder(t *testing.T) {
	res, err := http.Get(ts.URL + "/user/")
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
	id := []string{"@-1", "@-2", "@-3", "@@@"}
	for _, v := range id {
		res, err := http.Get(ts.URL + "/user/" + v)
		if err != nil {
			t.Fatal(err)
		}

		decoder := json.NewDecoder(res.Body)
		var myErr errors.Error
		err = decoder.Decode(&myErr)
		if err != nil {
			t.Fatal("Expected error in the body")
		}
		fmt.Println(v)
		if myErr.Code != errors.UserNotFound {
			t.Fatalf("Wrong err code: expected %v, got %v", errors.UserNotFound, myErr.Code)
		}
	}
}
