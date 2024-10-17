package render

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./dist/template/global/head.gohtml",
		"./dist/template/create-user.gohtml",
		"./dist/template/global/scripts.gohtml",
		"./dist/template/scripts/create-user.gohtml",
	)
	if err != nil {
		panic(err)
	}

	studentID := mux.Vars(r)["id"]
	if studentID == "" {
		http.Error(w, "Student ID is empty", http.StatusBadRequest)
		return
	}

	err = t.ExecuteTemplate(w, "create-user", struct {
		StudentID string
	}{
		StudentID: studentID,
	})
	if err != nil {
		panic(err)
	}
}
