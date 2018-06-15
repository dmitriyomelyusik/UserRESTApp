// Package controller has controller methods to operate with specific database.
// The interface of the database is wroten here.
// Then you can use that methods in your application handler.
package controller

import (
	"fmt"
	"strings"

	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

// Database is an interface that used in User controller to work with database data
type Database interface {
	GetUsers() ([]entity.User, error)
	GetUser(string) (entity.User, error)
	CreateUser(entity.User) (string, error)
	UpdateUser(entity.User) error
	DeleteUser(string) error
	UpdateUserFields(map[string]interface{}, string) error
}

// User controls database methods
type User struct {
	DB Database
}

// GetUsers is controlled method to get all users from database
func (ctl User) GetUsers() ([]entity.User, error) {
	return ctl.DB.GetUsers()
}

// GetUser is controlled method to get user with specific id from database
func (ctl User) GetUser(id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.Error{Code: errors.UserNotFound, Message: "invalid id"}
	}
	return ctl.DB.GetUser(id)
}

// CreateUser is controlled method to post new user
func (ctl User) CreateUser(user entity.User) (string, error) {
	return ctl.DB.CreateUser(user)
}

// UpdateUser is controlled method to change user info by its id
func (ctl User) UpdateUser(user entity.User) error {
	if _, err := ctl.DB.GetUser(user.ID); err != nil {
		return err
	}
	return ctl.DB.UpdateUser(user)
}

// DeleteUser is controlled method to delete user by its id
func (ctl User) DeleteUser(id string) error {
	if _, err := ctl.DB.GetUser(id); err != nil {
		return err
	}
	return ctl.DB.DeleteUser(id)
}

// UpdateUserFields is controlled method to patch user by its id with a set of instructions
func (ctl User) UpdateUserFields(changes map[string]interface{}, id string) error {
	if _, err := ctl.DB.GetUser(id); err != nil {
		return err
	}
	if len(changes) == 0 {
		return errors.Error{Code: errors.UserFieldsError, Message: "no data to change."}
	}
	var vrongFields []string
	var isOK = true
	for key := range changes {
		if key == "name" || key == "email" || key == "info" {
			continue
		}
		isOK = false
		vrongFields = append(vrongFields, key)
	}
	if !isOK {
		return errors.Error{Code: errors.UserFieldsError, Message: fmt.Sprintf("invalid user fields: %v.", strings.Join(vrongFields, ", "))}
	}
	return ctl.DB.UpdateUserFields(changes, id)
}
