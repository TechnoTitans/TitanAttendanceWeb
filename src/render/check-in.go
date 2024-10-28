package render

import (
	"TitanAttendance/src/auth"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func CheckIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/check-in.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/scripts/check-in.gohtml",
	)
	if err != nil {
		panic(err)
	}

	cookie, err := r.Cookie("TitanAttendancePin")
	if err != nil && err.Error() != "http: named cookie not present" {
		log.Error().Err(err).Msg("error getting cookie")
	}

	var userAccess = auth.PlainUser()
	if cookie != nil {
		userAccess, err = auth.CheckWithCookie(*cookie)
		if err != nil {
			log.Error().Err(err).Msg("error checking cookie")
		}
	}

	err = t.ExecuteTemplate(w, "check-in", struct {
		IsAdmin bool
	}{
		IsAdmin: userAccess.IsAdmin(),
	})
	if err != nil {
		panic(err)
	}
}
