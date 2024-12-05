package utils

import (
	"github.com/rs/zerolog/log"
	"os"
)

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	dbName        = os.Getenv("DB_NAME")
	domain        = os.Getenv("DOMAIN")
)

func init() {
	if adminPassword == "" {
		log.Fatal().Msg("ADMIN_PASSWORD environment variable not set.")
	}

	if dbName == "" {
		log.Fatal().Msg("DB_NAME environment variable not set.")
	}

	if domain == "" {
		log.Fatal().Msg("DOMAIN environment variable not set.")
	}
}

func GetAdminPassword() string {
	return adminPassword
}

func GetDBName() string {
	return dbName
}

func GetDomain() string {
	return domain
}
