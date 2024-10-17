package render

import (
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

	err = t.ExecuteTemplate(w, "check-in", struct{}{})
	if err != nil {
		panic(err)
	}
}
