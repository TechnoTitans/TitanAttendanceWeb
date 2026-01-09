package utils

import (
	"os"

	"github.com/rs/zerolog/log"
)

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	dbURL         = os.Getenv("DB_URL")
)

func init() {
	if adminPassword == "" {
		log.Fatal().Msg("ADMIN_PASSWORD environment variable not set.")
	}

	if dbURL == "" {
		log.Fatal().Msg("DB_URL environment variable not set.")
	}
}

func GetAdminPassword() string {
	return adminPassword
}

func GetDBURL() string {
	return dbURL
}
