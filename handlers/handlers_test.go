package handlers_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"."
	"../postgres"

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
	r = mux.NewRouter()
	r.HandleFunc("/user", handlers.UserHandler(p)).Methods("GET")
	r.HandleFunc("/user/{id}", handlers.UserIDHandler(p)).Methods("GET")
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

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	data := strings.Split(string(got), "\n")

	for i, v := range data {
		if i == len(data)-1 {
			break
		}
		str := strings.Fields(v)
		if len(str) < 4 {
			t.Fatal(str)
		}
	}
}

func TestUserIDHandler(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	for i := 1; i <= 50; i++ {
		res, err := http.Get(ts.URL + "/user/" + strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}

		got, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(got) != "" {
			str := strings.Fields(string(got))
			if str[0] != "ID"+strconv.Itoa(i) {
				t.Fatalf("Wrong id: expected ID%v, got %v", i, str[0])
			}
		}
	}

}
