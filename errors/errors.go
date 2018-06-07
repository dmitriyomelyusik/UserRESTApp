package errors

import "fmt"

//Error is custom error struct for convenient handlers work
type Error struct {
	Code    ErrCode
	Message string
	Info    interface{}
}

//ErrCode is code of Error that specify its
type ErrCode string

//All possible ErrCodes
const (
	UserNotFound       ErrCode = "UserNotFound"
	UsersNotFound      ErrCode = "UsersNotFound"
	UnmarshalError     ErrCode = "UnmarshallError"
	DatabaseQueryError ErrCode = "DatabaseQueryError"
	UserExists         ErrCode = "UserExists"
	NotFound           ErrCode = "NotFound"
)

func (e Error) Error() string {
	return fmt.Sprintf("%v\n%v\n%v", e.Code, e.Message, e.Info)
}
