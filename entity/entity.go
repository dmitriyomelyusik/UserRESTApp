// Package entity has entities (only one) that is used in application.
package entity

import (
	"fmt"
)

// User is type that contains user data
type User struct {
	ID    string      `json:"id,omitempty"`
	Email string      `json:"email"`
	Name  string      `json:"name"`
	Info  interface{} `json:"info,omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("ID%v %v %v \t %v", u.ID, u.Name, u.Email, u.Info)
}
