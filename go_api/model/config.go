package model

import (
	"os"
)

type AppConfig struct {
	Env                string
	Port               string
	UsersTableName     string
	FavoritesTableName string
	AwsEndpoint        string
	JwtSecretKey       string
}

// NewAppConfig creates a new config using environment variables
func NewAppConfig() AppConfig {
	return AppConfig{
		Env:                DefaultEnv("ENV", "local"),
		Port:               DefaultEnv("PORT", "8000"),
		UsersTableName:     DefaultEnv("USERS_TABLE_NAME", "the-drink-almanac-users"),
		FavoritesTableName: DefaultEnv("FAVORITES_TABLE_NAME", "the-drink-almanac-favorites"),
		AwsEndpoint:        os.Getenv("AWS_ENDPOINT"),
		JwtSecretKey:       os.Getenv("JWT_SECRET_KEY"),
	}
}

// DefaultEnv takes the name of the environment variable and a default value;
// if the environment variable wasn't found, then the default value is returned;
//
// Note: if the environment variable exists but just contains an empty string,
// the empty string will be returned
func DefaultEnv(envVarName, defaultValue string) string {
	envValue, ok := os.LookupEnv(envVarName)
	if !ok {
		envValue = defaultValue
	}
	return envValue
}
