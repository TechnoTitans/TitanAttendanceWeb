package main

import (
	"TitanAttendance/src/api"
	"TitanAttendance/src/datastore"
	"TitanAttendance/src/downloads"
	"TitanAttendance/src/middleware"
	"TitanAttendance/src/render"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate npm run build
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02 Jan 3:04:05 PM MST",
	})

	datastore.Connect(3 * time.Second)
	defer datastore.Disconnect()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", render.Login)
	router.NotFoundHandler = http.HandlerFunc(render.Error404)

	authenticatedRoute := router.NewRoute().Subrouter()
	authenticatedRoute.Use(middleware.Authenticate)
	authenticatedRoute.HandleFunc("/", render.CheckIn)
	authenticatedRoute.HandleFunc("/create-user/{id}", render.CreateUser)
	authenticatedRoute.HandleFunc("/qr", render.QRCode)

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/login", api.LogIn).Methods("POST")

	authenticatedApiRoute := apiRoute.NewRoute().Subrouter()
	authenticatedApiRoute.Use(middleware.Authenticate)
	authenticatedApiRoute.HandleFunc("/check-in", api.CheckIn).Methods("POST")
	authenticatedApiRoute.HandleFunc("/create-user", api.CreateUser).Methods("POST")
	authenticatedApiRoute.HandleFunc("/upload", api.Upload)
	authenticatedApiRoute.HandleFunc("/logout", api.LogOut).Methods("POST")

	downloadsRoute := router.PathPrefix("/downloads").Subrouter()
	downloadsRoute.HandleFunc("/export-database", downloads.ExportDatabase)

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Use(middleware.Sanitize)
	filesRoute.PathPrefix("/assets").Handler(
		http.StripPrefix("/files/assets", http.FileServer(http.Dir("public/assets"))),
	)

	log.Info().Msgf("Starting server.")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server.")
		return
	}
}
