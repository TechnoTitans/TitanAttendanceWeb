package api

import (
	"TitanAttendance/src/auth"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func CheckPin(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing body")
		}
	}(body)

	var req auth.Authentication
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

	_, err = req.CheckPin()
	if err != nil {
		if err.Error() == "invalid pin" {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}

		log.Error().Err(err).Msg("error checking pin")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	auth.SavePinCookie(w, req)
	w.WriteHeader(http.StatusOK)
}
