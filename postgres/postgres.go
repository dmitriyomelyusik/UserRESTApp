// Package postgres has a set of methods to operate with PostgreSQL database
package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

// Postgres is a PostgreSQL database
type Postgres struct {
	DB *sql.DB
}

// GetUsers returns all users in DB
func (p Postgres) GetUsers() ([]entity.User, error) {
	var users []entity.User
	rows, err := p.DB.Query("SELECT * FROM users")
	if err != nil {
		return users, errors.Error{Code: errors.DatabaseQueryError, Message: "invalid query to database"}
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, errors.Error{Code: errors.UsersNotFound, Message: "database didn't return users"}
		}
		if len(rawInfo) != 0 {
			err = json.Unmarshal(rawInfo, &u.Info)
			if err != nil {
				return users, errors.Error{Code: errors.UnmarshalError, Message: "postgres unmarshal error"}
			}
		}
		users = append(users, u)
	}
	return users, err
}

// GetUser returns user from DB with selected ID
func (p Postgres) GetUser(id string) (entity.User, error) {
	var u entity.User
	row := p.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var rawInfo []byte
	err := row.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)

	if err != nil {
		return u, errors.Error{Code: errors.UserNotFound, Message: fmt.Sprintf("cannot find user in database with id %v", id)}
	}
	if len(rawInfo) != 0 {
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return u, errors.Error{Code: errors.UnmarshalError, Message: "postgres unmarshal error"}
		}
	}
	return u, err
}

// UpdateUser eddits user by his id
func (p Postgres) UpdateUser(user entity.User) error {
	rawInfo, err := json.Marshal(user.Info)
	if err != nil {
		return errors.Error{Code: errors.UnmarshalError, Message: "postgres unmarshal error in UpdateUser method"}
	}
	res, err := p.DB.Exec("UPDATE users SET email=$1, name=$2, info=$3 WHERE id=$4", user.Email, user.Name, rawInfo, user.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: " + user.ID}
	}
	return nil
}

// CreateUser adds new user in database
func (p Postgres) CreateUser(user entity.User) (id string, err error) {
	id = generateID()
	rawInfo, err := json.Marshal(user.Info)
	if err != nil {
		return id, errors.Error{Code: errors.UnmarshalError, Message: "postgres unmarshal error in CreateUser method"}
	}
	res, err := p.DB.Exec("INSERT INTO users (id, email, name, info) VALUES ($1, $2, $3, $4)", id, user.Email, user.Name, rawInfo)
	if err != nil {
		return id, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return id, err
	}
	if n != 1 {
		return id, errors.Error{Code: errors.UnexpectedError, Message: "user is not created: " + id}
	}
	return id, nil
}

// DeleteUser deletes user by id
func (p Postgres) DeleteUser(id string) error {
	res, err := p.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: " + id}
	}
	return nil
}

// UpdateUserFields pathes modified user fields
func (p Postgres) UpdateUserFields(changes map[string]interface{}, id string) error {
	query := []string{"UPDATE users SET "}
	var fields []string
	for key, value := range changes {
		fields = append(fields, fmt.Sprintf("%v='%v'", key, value))
	}
	query = append(query, strings.Join(fields, ", "), fmt.Sprintf(" WHERE id='%v'", id))
	res, err := p.DB.Exec(strings.Join(query, ""))
	if err != nil {
		return errors.Error{Code: errors.DatabaseQueryError, Message: "bad query for database: " + strings.Join(query, "")}
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: " + id}
	}
	return nil
}

// NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

// generateID generates random id from letters and numbers
func generateID() string {
	signs := []rune("1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	id := make([]rune, 10)
	for i := range id {
		id[i] = signs[rnd.Intn(len(signs))]
	}
	return string(id)
}
