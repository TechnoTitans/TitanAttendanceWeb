package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

var users []User
var nameCaps = cases.Title(language.English, cases.Compact)

func AddNewStudent(user User) error {
	err := user.IsValid()
	if err != nil {
		return err
	}

	if user.IDExists() {
		return errors.New("student ID already exists")
	}

	user.Name = strings.Join(strings.Fields(user.Name), " ")
	user.Name = nameCaps.String(user.Name)

	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = conn.Database(utils.DBName).Collection("students").InsertOne(ctx, user)
	if err == nil {
		users = append(users, user)
	}

	log.Info().Msgf("Added %s as a student.", user.Name)
	return err
}

func GetStudents() []User {
	if len(users) > 0 {
		return users
	}

	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := conn.Database(utils.DBName).Collection("students").Find(ctx, map[string]interface{}{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get students.")
		return nil
	}
	defer func(cur *mongo.Cursor) {
		err = cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to close cursor.")
		}
	}(cur)

	for cur.Next(ctx) {
		var user User
		err = cur.Decode(&user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode user.")
			return nil
		}
		users = append(users, user)
	}

	return users
}

func ClearAllStudents() error {
	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Database(utils.DBName).Collection("students").DeleteMany(ctx, nil)
	if err != nil {
		return err
	}

	users = []User{}
	return nil
}
