package main

import (
	"TitanAttendance/src/api"
	"TitanAttendance/src/database"
	"TitanAttendance/src/middleware"
	"TitanAttendance/src/render"
	"TitanAttendance/src/users"
	"TitanAttendance/src/utils"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02 Jan 3:04:05 PM MST",
	})

	db := database.GetFireDB()
	err := db.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Firebase.")
	}
	defer db.Disconnect()

	if len(os.Args) >= 2 {
		var method utils.CSVUploadMethod
		method, err = utils.AskForNewCSVMethod()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get CSV upload method.")
			<-time.After(5 * time.Second)
			return
		}
		users.AddStudentsFromCSV(os.Args[1], method)
		<-time.After(5 * time.Second)
		return
	}

	users.GetStudents()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/pin", render.CheckPin)
	router.NotFoundHandler = http.HandlerFunc(render.Error404)

	authenticatedRoute := router.NewRoute().Subrouter()
	authenticatedRoute.Use(middleware.Authenticate)
	authenticatedRoute.HandleFunc("/", render.CheckIn)
	authenticatedRoute.HandleFunc("/create-user/{id}", render.CreateUser)
	authenticatedRoute.HandleFunc("/qr", render.QRCode)

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/check-pin", api.CheckPin).Methods("POST", "OPTIONS")

	authenticatedApiRoute := apiRoute.NewRoute().Subrouter()
	authenticatedApiRoute.Use(middleware.Authenticate)
	authenticatedApiRoute.HandleFunc("/check-in", api.CheckIn).Methods("POST", "OPTIONS")
	authenticatedApiRoute.HandleFunc("/create-user", api.CreateUser).Methods("POST", "OPTIONS")

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Use(middleware.Sanitize)
	filesRoute.PathPrefix("/assets").Handler(
		http.StripPrefix("/files/assets", http.FileServer(http.Dir("public/assets"))),
	)

	log.Info().Msgf("Starting server.")
	err = http.ListenAndServe(utils.GetPort(), router)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server.")
		return
	}
}
