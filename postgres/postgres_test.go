package postgres_test

import (
	"log"
	"os"
	"testing"

	"github.com/UserRESTApp/entity"

	. "github.com/UserRESTApp/postgres"
	_ "github.com/lib/pq"
)

var (
	p  *Postgres
	ID []string
)

func TestMain(m *testing.M) {
	var err error
	p, err = NewDB("user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestUsers(t *testing.T) {
	u, err := p.GetUsers()
	if err != nil {
		t.Fatalf("Cannot select any of users: %v", err)
	}
	for _, v := range u {
		if v.ID == "" {
			t.Fatalf("ID must be not nil")
		}
		ID = append(ID, v.ID)
	}
}

func TestUserByID(t *testing.T) {
	for _, id := range ID {
		_, err := p.GetUserByID(id)
		if err != nil {
			t.Fatalf("Cannot get existing user: %v", err)
		}
	}
}

func TestPostUser(t *testing.T) {
	var test = entity.User{ID: "test", Name: "test", Email: "test", Info: "\"test\""}
	if err := p.PostUser(test); err != nil {
		t.Fatal(err)
	}
	test.Info = "test"
	new, err := p.GetUserByID("test")
	if err != nil {
		t.Fatal(err)
	}
	if new != test {
		t.Fatal("Cannot post user into database")
	}
}

func TestPutUser(t *testing.T) {
	var test = entity.User{ID: "test", Name: "test2", Email: "test2", Info: "\"test2\""}
	if err := p.PutUserByID(test); err != nil {
		t.Fatal(err)
	}
	test.Info = "test2"
	new, err := p.GetUserByID("test")
	if err != nil {
		t.Fatal(err)
	}
	if new != test {
		t.Fatal("Cannot put user into database")
	}
}

func TestPatchUser(t *testing.T) {
	m := make(map[string]interface{})
	m["Name"] = "test3"
	var test = entity.User{ID: "test", Name: "test3", Email: "test2", Info: "test2"}
	if err := p.PatchUserByID(m, "test"); err != nil {
		t.Fatal(err)
	}
	new, err := p.GetUserByID("test")
	if err != nil {
		t.Fatal(err)
	}
	if new != test {
		t.Fatal("Cannot patch user into database")
	}
}

func TestDeleteUser(t *testing.T) {
	if err := p.DeleteUserByID("test"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.GetUserByID("test"); err == nil {
		t.Fatal("User is not deleted")
	}
}
