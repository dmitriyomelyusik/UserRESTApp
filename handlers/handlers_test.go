package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func TestUserHanlder(t *testing.T) {
	res, err := http.Get(ts.URL + "/user/")
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	var users []entity.User
	err = decoder.Decode(&users)
	if err != nil {
		t.Fatal("Bad user perfomance")
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
		if myErr.Code != errors.UserNotFound {
			t.Fatalf("Wrong err code: expected %v, got %v", errors.UserNotFound, myErr.Code)
		}
	}
}

func TestPostUser(t *testing.T) {
	r := bytes.NewReader([]byte(`{"Name":"test","Email":"test","ID":"test","Info":"\"test\""}`))
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/user/", r)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("Cannot post user")
	}
	res, err = http.Get(ts.URL + "/user/test")
	if err != nil {
		t.Fatal(err)
	}
	var u entity.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "test" || u.Email != "test" || u.ID != "test" || u.Info != "test" {
		t.Fatalf("Cannot post user %s", u)
	}
}

func TestPutUser(t *testing.T) {
	r := bytes.NewReader([]byte(`{"Name":"test2","Email":"test2","ID":"test","Info":"\"test2\""}`))
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/user/test", r)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("Cannot put user")
	}
	res, err = http.Get(ts.URL + "/user/test")
	if err != nil {
		t.Fatal(err)
	}
	var u entity.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "test2" || u.Email != "test2" || u.ID != "test" || u.Info != "test2" {
		t.Fatalf("Cannot put user %s", u)
	}
}

func TestPatchUser(t *testing.T) {
	r := bytes.NewReader([]byte(`{"Name":"test3"}`))
	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/user/test", r)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("Cannot patch user")
	}
	if err != nil {
		t.Fatal(err)
	}
	res, err = http.Get(ts.URL + "/user/test")
	if err != nil {
		t.Fatal(err)
	}
	var u entity.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "test3" || u.Email != "test2" || u.ID != "test" || u.Info != "test2" {
		t.Fatalf("Cannot post user %s", u)
	}
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("Cannot delete user")
	}
	res, err = http.Get(ts.URL + "/user/test")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatal("Test user is not deleted")
	}
}
