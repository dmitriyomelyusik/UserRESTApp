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
