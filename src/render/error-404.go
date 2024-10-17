package render

import (
	"html/template"
	"net/http"
)

func Error404(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/error-404.gohtml",
		"./dist/template/global/scripts.gohtml")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNotFound)

	err = t.ExecuteTemplate(w, "error-404", struct{}{})
	if err != nil {
		panic(err)
	}
}
