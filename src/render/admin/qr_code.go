package admin

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/utils"
	"html/template"
	"net/http"
)

func QRCode(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/qr-code.gohtml",
		"./dist/template/global/scripts.gohtml",
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

	newPin := auth.CreateUserPin()

	err = t.ExecuteTemplate(w, "qrcode", struct {
		QRCode string
		Pin    string
	}{
		QRCode: utils.CreateQRCode(&newPin),
		Pin:    newPin,
	})
	if err != nil {
		panic(err)
	}
}
