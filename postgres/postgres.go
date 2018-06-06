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
func (p Postgres) GetUsers() ([]entity.User, errors.Error) {
	var users []entity.User
	rows, err := p.DB.Query("SELECT * FROM users")
	if err != nil {
		return users, errors.Error{Code: errors.DatabaseQueryError, Message: "Invalid query to database"}
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, errors.Error{Code: errors.UsersNotFound, Message: "Database didn't return users"}
		}
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return users, errors.Error{Code: errors.UnmarshalError, Message: "Something was wrong in unmarshalling process"}
		}
		users = append(users, u)
	}
	return users, errors.Error{}
}

//GetUserByID returns user from DB with selected ID
func (p Postgres) GetUserByID(id string) (entity.User, errors.Error) {
	var u entity.User
	row := p.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var rawInfo []byte
	err := row.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)

	if err != nil {
		return u, errors.Error{Code: errors.UserNotFound, Message: "Cannot find user in database with that id"}
	}
	err = json.Unmarshal(rawInfo, &u.Info)
	if err != nil {
		err = errors.Error{Code: errors.UnmarshalError, Message: "Something was wrong in unmarshalling process"}
	}
	return u, errors.Error{}
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
