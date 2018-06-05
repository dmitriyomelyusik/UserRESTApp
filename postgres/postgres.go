package postgres

import (
	"database/sql"
	"encoding/json"

	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

//Postgres is a PostgreSQL database
type Postgres struct {
	DB *sql.DB
}

//GetUsers returns all users in DB
func (p Postgres) GetUsers() ([]entity.User, error) {
	var users []entity.User
	rows, err := p.DB.Query("SELECT * FROM users")
	if err != nil {
		return users, errors.ErrDatabaseQuery
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, errors.ErrUsersNotFound
		}
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return users, errors.ErrUnmarshal
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
		return u, errors.ErrUserNotFound
	}
	err = json.Unmarshal(rawInfo, &u.Info)
	if err != nil {
		err = errors.ErrUnmarshal
	}
	return u, err
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
