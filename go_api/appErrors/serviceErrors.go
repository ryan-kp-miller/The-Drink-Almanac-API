package appErrors

type UserAlreadyExistsError struct {
	message string
}

func (e UserAlreadyExistsError) Error() string {
	return e.message
}

func NewUserAlreadyExistsError(message string) UserAlreadyExistsError {
	return UserAlreadyExistsError{message: message}
}
