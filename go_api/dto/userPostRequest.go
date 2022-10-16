package dto

import "fmt"

type UserPostRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u UserPostRequest) ValidateRequest() error {
	if u.Username == "" {
		return fmt.Errorf("no username provided")
	}
	if u.Password == "" {
		return fmt.Errorf("no password provided")
	}
	return nil
}
