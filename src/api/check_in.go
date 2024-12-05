package api

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/users"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func CheckIn(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing body")
		}
	}(body)

	var user users.User
	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("error decoding request")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	if user.ID == auth.GetPin() {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("your id cant be the same as the pin"))
		return
	}

	err = user.CheckIn()
	if err != nil {
		log.Warn().Err(err).Msg("error checking in")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
