package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"time"
)

const authUrl = "mongodb://127.0.0.1:27017/"

var client *mongo.Client

func Connect(timeout time.Duration) {
	var err error

	var clientOptions = options.Client().
		ApplyURI(fmt.Sprintf(authUrl)).
		SetMinPoolSize(15).
		SetWriteConcern(writeconcern.Majority()).
		SetReadConcern(readconcern.Majority())

	client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to establish connection to MongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to establish connection to MongoDB")
	} else {
		log.Info().Msg("Established connection to MongoDB")
	}
}

func Disconnect() {
	if client != nil {
		_ = client.Disconnect(context.Background())
	}
}

func GetConn() *mongo.Client {
	return client
}
