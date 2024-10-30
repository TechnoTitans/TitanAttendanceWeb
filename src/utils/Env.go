package utils

import "os"

const (
	port   = ":8081"
	dbName = "titanattendance"
)

func GetPort() string {
	return port
}

func GetDBName() string {
	return dbName
}

func GetDomain() string {
	return "https://spira1.com"
}

var adminPassword = os.Getenv("ADMIN_PASSWORD")

func GetAdminPassword() string {
	return adminPassword
}
