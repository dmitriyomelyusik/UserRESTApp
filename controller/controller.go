package controller

import (
	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
)

//Database is an interface that used in User controller to work with database data
type Database interface {
	GetUsers() ([]entity.User, error)
	GetUserByID(string) (entity.User, error)
}

//User is controlled
type User struct {
	DB Database
}

//GetUsers is controlled method to get all users from database
func (ctl User) GetUsers() ([]entity.User, error) {
	users, err := ctl.DB.GetUsers()
	checkError(&err)
	return users, err
}

//GetUserByID is controlled method to get user with specific id from database
func (ctl User) GetUserByID(id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.ErrNotFound
	}
	user, err := ctl.DB.GetUserByID(id)
	checkError(&err)
	return user, err
}

func checkError(err *error) {
	switch *err {
	case errors.ErrDatabaseQuery:
		*err = errors.ErrInternalServer
	case errors.ErrUnmarshal:
		*err = errors.ErrInternalServer
	case errors.ErrUsersNotFound:
		*err = errors.ErrInternalServer
	case errors.ErrUserNotFound:
		*err = errors.ErrNotFound
	}
}
