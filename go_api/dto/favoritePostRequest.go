package dto

import "fmt"

type FavoritePostRequest struct {
	UserId  string `json:"user_id"`
	DrinkId string `json:"drink_id"`
}

func (f FavoritePostRequest) ValidateRequest() error {
	if f.UserId == "" {
		return fmt.Errorf("no user id provided")
	}
	if f.DrinkId == "" {
		return fmt.Errorf("no drink id provided")
	}
	return nil
}
