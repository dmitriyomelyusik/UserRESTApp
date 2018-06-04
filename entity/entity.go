package entity

import (
	"fmt"
)

//User is type that contains user data
type User struct {
	ID    string
	Email string
	Name  string
	Info  interface{}
}

func (u User) String() string {
	return fmt.Sprintf("ID%v %v %v \t %v", u.ID, u.Name, u.Email, u.Info)
}

//Database is an interface that provide you access to operations with database
type Database interface {
	Users() ([]User, error)
	GetUserByID(id string) (User, error)
}
