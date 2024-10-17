package utils

import "os"

const (
	port = ":8081"
)

func GetPort() string {
	return port
}

// https://console.firebase.google.com/u/3/project/titanattendance-91b96/firestore/databases/-default-/data/~2Fmeetings
var (
	apiKey        = os.Getenv("API_KEY")
	authDomain    = os.Getenv("AUTH_DOMAIN")
	projectID     = os.Getenv("PROJECT_ID")
	storageBucket = os.Getenv("STORAGE_BUCKET")
)

func GetAPIKey() string {
	return apiKey
}

func GetAuthDomain() string {
	return authDomain
}

func GetProjectID() string {
	return projectID
}

func GetStorageBucket() string {
	return storageBucket
}

func GetDomain() string {
	return "https://spira1.com"
}
