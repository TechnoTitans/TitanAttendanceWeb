package database

import (
	"TitanAttendance/src/utils"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type FireDB struct {
	*firestore.Client
}

var fireDB FireDB

func (db *FireDB) Connect() error {
	ctx := context.Background()
	opt := option.WithAPIKey(utils.GetAPIKey())
	config := &firebase.Config{
		DatabaseURL:   utils.GetAuthDomain(),
		ProjectID:     utils.GetProjectID(),
		StorageBucket: utils.GetStorageBucket(),
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	db.Client = client
	log.Info().Msgf("Connected to Firebase.")
	return nil
}

func (db *FireDB) Disconnect() {
	if db.Client != nil {
		err := db.Client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to disconnect from Firebase.")
		} else {
			log.Info().Msg("Disconnected from Firebase.")
		}
	}
}

func GetFireDB() *FireDB {
	return &fireDB
}
