package errors

import "fmt"

//Error bla-bla
type Error struct {
	Code    ErrCode
	Message string
	Info    interface{}
}

//ErrCode bla-bla
type ErrCode string

//All possible ErrCodes
const (
	NotFound            ErrCode = "NotFound"
	InternalServerError ErrCode = "InternalServerError"
	UserNotFound        ErrCode = "UserNotFound"
	UsersNotFound       ErrCode = "UsersNotFound"
	UnmarshalError      ErrCode = "UnmarshallError"
	DatabaseQueryError  ErrCode = "DatabaseQueryError"
)

func (e Error) Error() string {
	return fmt.Sprintf("%v\n%v\n%v", e.Code, e.Message, e.Info)
}
