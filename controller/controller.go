package controller

import (
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

//Database is an interface that used in User controller to work with database data
type Database interface {
	GetUsers() ([]entity.User, error)
	GetUserByID(string) (entity.User, error)
	PostUser(entity.User) error
	PutUserByID(entity.User) error
	DeleteUserByID(string) error
}

//User controls database methods
type User struct {
	DB Database
}

//GetUsers is controlled method to get all users from database
func (ctl User) GetUsers() ([]entity.User, error) {
	return ctl.DB.GetUsers()
}

//GetUserByID is controlled method to get user with specific id from database
func (ctl User) GetUserByID(id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.Error{Code: errors.UserNotFound, Message: "Invalid id"}
	}
	return ctl.DB.GetUserByID(id)
}

//PostUser is controlled method to post new user
func (ctl User) PostUser(user entity.User) error {
	u, err := ctl.DB.GetUserByID(user.ID)
	if u != (entity.User{}) {
		return errors.Error{Code: errors.UserExists, Message: "Cannot add new user: user with that id exists."}
	}
	myErr, ok := err.(errors.Error)
	if !ok {
		return err
	}
	if myErr.Code != errors.UserNotFound {
		return myErr
	}
	return ctl.DB.PostUser(user)
}

//PutUserByID is controlled method to change user info by its id
func (ctl User) PutUserByID(user entity.User, id string) error {
	if id != user.ID {
		return errors.Error{Code: errors.NotFound, Message: "User has invalid id."}
	}
	u, err := ctl.DB.GetUserByID(user.ID)
	if u == (entity.User{}) {
		return ctl.PostUser(user)
	}
	if err != nil {
		return err
	}
	return ctl.DB.PutUserByID(user)
}

//DeleteUserByID is controlled method to delete user by its id
func (ctl User) DeleteUserByID(id string) error {
	_, err := ctl.DB.GetUserByID(id)
	if err != nil {
		return err
	}
	return ctl.DB.DeleteUserByID(id)
}
