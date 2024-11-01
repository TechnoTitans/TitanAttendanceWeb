package utils

import (
	"github.com/rs/zerolog/log"
	"os"
)

const (
	DBName = "titanattendance"
	Domain = "https://spira1.com"
)

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

func init() {
	if adminPassword == "" {
		log.Fatal().Msg("ADMIN_PASSWORD environment variable not set.")
	}
}

func GetAdminPassword() string {
	return adminPassword
}
