package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

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
		return users, errors.Error{Code: errors.DatabaseQueryError, Message: "Invalid query to database"}
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, errors.Error{Code: errors.UsersNotFound, Message: "Database didn't return users"}
		}
		if len(rawInfo) != 0 {
			err = json.Unmarshal(rawInfo, &u.Info)
			if err != nil {
				return users, errors.Error{Code: errors.UnmarshalError, Message: "Something was wrong in unmarshalling process"}
			}
		}
		users = append(users, u)
	}
	return users, err
}

//GetUserByID returns user from DB with selected ID
func (p Postgres) GetUserByID(id string) (entity.User, error) {
	var u entity.User
	row := p.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var rawInfo []byte
	err := row.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)

	if err != nil {
		return u, errors.Error{Code: errors.UserNotFound, Message: fmt.Sprintf("Cannot find user in database with id %v", id)}
	}
	if len(rawInfo) != 0 {
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return u, errors.Error{Code: errors.UnmarshalError, Message: "Something was wrong in unmarshalling process"}
		}
	}
	return u, err
}

//PutUserByID eddits user by his id
func (p Postgres) PutUserByID(user entity.User) error {
	_, err := p.DB.Exec("UPDATE users SET email=$1, name=$2, info=$3 WHERE id=$4", user.Email, user.Name, user.Info, user.ID)
	return err
}

//PostUser adds new user in database
func (p Postgres) PostUser(user entity.User) error {
	_, err := p.DB.Exec("INSERT INTO users (id, email, name, info) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.Name, user.Info)
	return err
}

//DeleteUserByID deletes user by id
func (p Postgres) DeleteUserByID(id string) error {
	_, err := p.DB.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

//PatchUserByID pathes modified user fields
func (p Postgres) PatchUserByID(changes map[string]interface{}, id string) error {
	var query = "UPDATE users SET "
	for key, value := range changes {
		query += fmt.Sprintf("%v='%v', ", key, value)
	}
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id='%v'", id)
	_, err := p.DB.Exec(query)
	return err
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
