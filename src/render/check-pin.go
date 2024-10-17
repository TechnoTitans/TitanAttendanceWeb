package render

import (
	"TitanAttendance/src/auth"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func CheckPin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/pin.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/scripts/pin.gohtml",
	)
	if err != nil {
		panic(err)
	}

	pin := r.URL.Query().Get("code")
	if pin != "" {
		userAuth := auth.Authentication{Pin: pin}
		_, err = userAuth.CheckPin()
		if err != nil {
			if err.Error() != "invalid pin" {
				log.Error().Err(err).Msg("error checking pin")
				fmt.Println(2)
			}
		} else {
			auth.SavePinCookie(w, userAuth)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	cookie, err := r.Cookie("TitanAttendancePin")
	if err != nil && err.Error() != "http: named cookie not present" {
		log.Error().Err(err).Msg("error getting cookie")
	}

	if cookie != nil {
		userAccess, err := auth.CheckWithCookie(*cookie)
		if err == nil && userAccess.IsAllowed() {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	err = t.ExecuteTemplate(w, "pin", struct{}{})
	if err != nil {
		panic(err)
	}
}
