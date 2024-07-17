package dto

import "fmt"

type FavoritePostRequest struct {
	DrinkId string `json:"drink_id"`
}

func (f FavoritePostRequest) ValidateRequest() error {
	if f.DrinkId == "" {
		return fmt.Errorf("no drink id provided")
	}
	return nil
}
