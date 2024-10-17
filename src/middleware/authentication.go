package middleware

import (
	"TitanAttendance/src/auth"
	"net/http"
)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TitanAttendancePin")
		if err != nil {
			http.Redirect(w, r, "/pin", http.StatusFound)
			return
		}

		_, err = auth.CheckWithCookie(*cookie)
		if err != nil {
			http.Redirect(w, r, "/pin", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}
