package postgres_test

import (
	"fmt"
	"os"
	"testing"

	"."
	_ "github.com/lib/pq"
)

var (
	p  *postgres.Postgres
	ID []string
)

func TestMain(m *testing.M) {
	var err error
	p, err = postgres.NewDB("host=127.0.0.1 user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		fmt.Printf("Cannot open database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestUsers(t *testing.T) {
	u, err := p.Users()
	if err != nil {
		t.Fatalf("Cannot select any user: %v", err)
	}
	for _, v := range u {
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
