package api

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/users"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response.")
		}
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

	err = users.ClearAllStudents()
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response.")
		}
		return
	}

	err = users.ClearAllMeetings()
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response.")
		}
		return
	}

	failedToAdd := 0
	for _, row := range rows {
		var user users.User

		user = users.User{
			ID:   row[idsInColumn],
			Name: row[1-idsInColumn],
		}

		err = users.AddNewStudent(user)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to add %s as a student.", user.Name)
			failedToAdd++
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("Added %d students. Failed to add %d.", len(rows)-failedToAdd, failedToAdd)))
	if err != nil {
		log.Error().Err(err).Msg("Failed to write response.")
		return
	}
}
