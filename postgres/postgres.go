package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/UserRESTApp/entity"
)

//Postgres is a PostgreSQL database
type Postgres struct {
	DB *sql.DB
}

const (
	userNotFoundError  = "User is not found"
	usersNotFoundError = "Users are not found"
	unmarshalError     = "Cannot unmarshal data"
	databaseQueryError = "Database query error"
)

//Users returns all users in DB
func (p Postgres) Users() ([]entity.User, error) {
	var users []entity.User
	rows, err := p.DB.Query("SELECT * FROM users")
	if err != nil {
		return users, errors.New(databaseQueryError)
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, errors.New(usersNotFoundError)
		}
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return users, errors.New(unmarshalError)
		}
		users = append(users, u)
	}
	return users, nil
}

//GetUserByID returns user from DB with selected ID
func (p Postgres) GetUserByID(id string) (entity.User, error) {
	var u entity.User
	row := p.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var rawInfo []byte
	err := row.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)

	if err != nil {
		return u, errors.New(userNotFoundError)
	}
	err = json.Unmarshal(rawInfo, &u.Info)
	return u, errors.New(unmarshalError)
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
