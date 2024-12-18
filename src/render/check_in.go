package render

import (
	"TitanAttendance/src/auth"
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
)

func CheckIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/global/nav.gohtml",
		"./dist/template/check-in.gohtml",
		"./dist/template/scripts/check-in.gohtml",
		"./dist/template/scripts/pines/toast.gohtml",
		"./dist/template/scripts/pines/new-user-modal.gohtml",
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
