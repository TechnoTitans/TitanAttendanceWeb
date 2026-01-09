package utils

import "time"

func GetCurrentDate() string {
	return time.Now().Format("01-02-2006")
}

func GetCurrentTime() string {
	return time.Now().Format("3:04:05 PM")
}
