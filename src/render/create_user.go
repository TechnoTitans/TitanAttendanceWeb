package render

import (
	"TitanAttendance/src/auth"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/global/nav.gohtml",
		"./dist/template/create-user.gohtml",
		"./dist/template/scripts/create-user.gohtml",
		"./dist/template/scripts/pines/toast.gohtml",
	)
	if err != nil {
		panic(err)
	}

	studentID := mux.Vars(r)["id"]
	if studentID == "" {
		http.Error(w, "Student ID is empty", http.StatusBadRequest)
		return
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

	err = t.ExecuteTemplate(w, "create-user", struct {
		StudentID string
		IsAdmin   bool
	}{
		StudentID: studentID,
		IsAdmin:   userAccess.IsAdmin(),
	})
	if err != nil {
		panic(err)
	}
}
