package appErrors

import "fmt"

type UserAlreadyExistsError struct {
	message string
}

func (e UserAlreadyExistsError) Error() string {
	return e.message
}

func NewUserAlreadyExistsError(username string) UserAlreadyExistsError {
	return UserAlreadyExistsError{message: fmt.Sprintf("a user already exists with the username '%s'", username)}
}

type FavoriteAlreadyExistsError struct {
	message string
}

func (e FavoriteAlreadyExistsError) Error() string {
	return e.message
}

func NewFavoriteAlreadyExistsError(message string) FavoriteAlreadyExistsError {
	return FavoriteAlreadyExistsError{message: message}
}

type UserNotFoundError struct {
	message string
}

func (e UserNotFoundError) Error() string {
	return e.message
}

func NewUserNotFoundError(username string) UserNotFoundError {
	return UserNotFoundError{message: fmt.Sprintf("no user exists with username '%s'", username)}
}

type InvalidAuthTokenError struct {
	message string
}

func (e InvalidAuthTokenError) Error() string {
	return e.message
}

func NewInvalidAuthTokenError(message string) InvalidAuthTokenError {
	return InvalidAuthTokenError{message: message}
}

type IncorrectPasswordError struct {
	message string
}

func (e IncorrectPasswordError) Error() string {
	return e.message
}

func NewIncorrectPasswordError(username string) IncorrectPasswordError {
	return IncorrectPasswordError{message: fmt.Sprintf("incorrect password for username '%s'", username)}
}
