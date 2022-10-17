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

type FavoriteAlreadyExistsError struct {
	message string
}

func (e FavoriteAlreadyExistsError) Error() string {
	return e.message
}

func NewFavoriteAlreadyExistsError(message string) FavoriteAlreadyExistsError {
	return FavoriteAlreadyExistsError{message: message}
}
