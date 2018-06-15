// Package errors has custom wroten Error for convenience.
package errors

import "fmt"

// Error is custom error struct for convenient handlers work
type Error struct {
	Code    ErrCode     `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}

// ErrCode is code of Error that specify its
type ErrCode string

// All possible ErrCodes
const (
	UserNotFound       ErrCode = "userNotFound"
	UsersNotFound      ErrCode = "usersNotFound"
	UnmarshalError     ErrCode = "unmarshallError"
	DatabaseQueryError ErrCode = "databaseQueryError"
	UserExists         ErrCode = "userExists"
	NotFound           ErrCode = "notFound"
	UserFieldsError    ErrCode = "userFieldsError"
	UnexpectedError    ErrCode = "unexpectedError"
)

func (e Error) Error() string {
	return fmt.Sprintf("%v\n%v\n%v", e.Code, e.Message, e.Info)
}
