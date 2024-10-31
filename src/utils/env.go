package utils

import "os"

const (
	DBName = "titanattendance"
	Domain = "https://spira1.com"
)

var (
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

func GetAdminPassword() string {
	return adminPassword
}
