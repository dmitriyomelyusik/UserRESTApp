package errors

import "errors"

//Errors that may occure in postgres methods
var (
	ErrUserNotFound  = errors.New("User is not found")
	ErrUsersNotFound = errors.New("Users are not found")
	ErrUnmarshal     = errors.New("Cannot unmarshal data")
	ErrDatabaseQuery = errors.New("Database query error")
)

//Errors that is handled at handlers
var (
	ErrInternalServer = errors.New("Internal server error")
	ErrNotFound       = errors.New("Not found error")
)
