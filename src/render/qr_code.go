package render

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/utils"
	"html/template"
	"net/http"
)

func QRCode(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/global/nav.gohtml",
		"./dist/template/qr-code.gohtml",
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

	newPin := auth.GetPin()

	err = t.ExecuteTemplate(w, "qrcode", struct {
		QRCode  string
		Pin     string
		IsAdmin bool
	}{
		QRCode:  utils.CreateQRCode(&newPin).Base64,
		Pin:     newPin,
		IsAdmin: userAccess.IsAdmin(),
	})
	if err != nil {
		panic(err)
	}
}
