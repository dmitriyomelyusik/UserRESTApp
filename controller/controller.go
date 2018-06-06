package controller

import (
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

//Database is an interface that used in User controller to work with database data
type Database interface {
	GetUsers() ([]entity.User, errors.Error)
	GetUserByID(string) (entity.User, errors.Error)
}

//User controls database methods
type User struct {
	DB Database
}

//GetUsers is controlled method to get all users from database
func (ctl User) GetUsers() ([]entity.User, errors.Error) {
	return ctl.DB.GetUsers()
}

//GetUserByID is controlled method to get user with specific id from database
func (ctl User) GetUserByID(id string) (entity.User, errors.Error) {
	if id == "" {
		return entity.User{}, errors.Error{Code: errors.NotFound, Message: "Invalid id"}
	}
	return ctl.DB.GetUserByID(id)
}
