package api

import (
	"TitanAttendance/src/auth"
	"net/http"
)

func LogOut(w http.ResponseWriter, _ *http.Request) {
	auth.ClearPinCookie(w)
	w.WriteHeader(http.StatusOK)
}
