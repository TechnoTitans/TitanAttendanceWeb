package render

import (
	"TitanAttendance/src/auth"
	"html/template"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/global/nav.gohtml",
		"./dist/template/scripts/upload.gohtml",
		"./dist/template/upload.gohtml",
	)
	if err != nil {
		panic(err)
	}

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

	err = t.ExecuteTemplate(w, "upload", struct {
		IsAdmin bool
	}{
		IsAdmin: userAccess.IsAdmin(),
	})
	if err != nil {
		panic(err)
	}
}
