package datastore

import (
	"TitanAttendance/src/utils"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var client *pgxpool.Pool

func Connect(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := pgxpool.New(ctx, utils.GetDBURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database.")
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to ping database. URL: %v", utils.GetDBURL())
	} else {
		log.Info().Msg("Successfully pinged database.")
	}
	_, err = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS students (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create students table.")
	}

	_, err = db.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS meetings (date TEXT PRIMARY KEY, absent JSONB NOT NULL, present JSONB NOT NULL)`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create meetings table.")
	}

	client = db
}

func Disconnect() {
	client.Close()
}

func GetClient() *pgxpool.Pool {
	return client
}
