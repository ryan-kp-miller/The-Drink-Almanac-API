package lambda

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"the-drink-almanac-api/service"
)

// authorizeUser extracts a userId from the Token header and returns the userId if they are authorized
func authorizeUser(headers map[string]string, authService service.AuthService) (string, error) {
	token := headers["Token"]
	if len(token) == 0 {
		return "", errors.New("the 'Token' header was not included in the request")
	}

	userId, err := authService.ValidateToken(token)
	if err != nil {
		return "", errors.New("the 'Token' header was invalid")
	}

	return userId, nil
}

func messageToResponseBody(message string) string {
	m := map[string]string{
		"message": message,
	}
	body, _ := jsoniter.MarshalToString(m)
	return body
}