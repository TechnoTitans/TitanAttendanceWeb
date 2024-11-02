package api

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/users"
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

func Upload(w http.ResponseWriter, r *http.Request) {
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

	body := r.Body
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing body")
		}
	}(body)

	reader := csv.NewReader(r.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read file.")
		return
	}

	idsInColumn := 0
	for _, row := range rows {
		for col, cell := range row {
			_, err = strconv.Atoi(cell)
			if err == nil {
				idsInColumn = col
				break
			}
		}
	}

	for _, row := range rows {
		var user users.User

		user = users.User{
			ID:   row[idsInColumn],
			Name: row[1-idsInColumn],
		}

		err = users.AddNewStudent(user)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to add %s as a student.", user.Name)
		}
	}

	w.WriteHeader(http.StatusOK)
}
