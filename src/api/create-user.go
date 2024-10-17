package api

import (
	"TitanAttendance/src/users"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing body")
		}
	}(body)

	var req users.User
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("error decoding request")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	err = users.AddNewStudent(req)
	if err != nil {
		log.Error().Err(err).Msg("error adding new student")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	err = req.CheckIn()
	if err != nil {
		log.Error().Err(err).Msg("error checking in")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
