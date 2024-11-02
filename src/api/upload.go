package api

import (
	"TitanAttendance/src/auth"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("TitanAttendancePin")
	if err != nil {
		return
	}

	userAccess, err := auth.CheckWithCookie(*cookie)
	if err != nil {
		return
	}

	if !userAccess.IsAdmin() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

}
