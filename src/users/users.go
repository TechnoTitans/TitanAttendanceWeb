package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"context"
	"encoding/csv"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"os"
	"strings"
	"time"
)

var users []User

func AddNewStudent(user User) error {
	err := user.IsValid()
	if err != nil {
		return err
	}

	if user.IDExists() {
		return errors.New("student ID already exists")
	}

	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.Name = strings.Join(strings.Fields(user.Name), " ")

	_, err = conn.Database(utils.DBName).Collection("students").InsertOne(ctx, user)
	if err == nil {
		users = append(users, user)
	}

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

func AddStudentsFromCSV(filePath string, method utils.CSVUploadMethod) {
	log.Info().Msgf("Adding students from %s", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Err(err).Msg("File does not exist.")
		return
	}

	open, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file.")
		return
	}
	defer func(open *os.File) {
		err = open.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close file.")
			return
		}
	}(open)

	reader := csv.NewReader(open)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read file.")
		return
	}

	if method.ClearAllPreviousIDs {
		err = ClearAllMeetings()
		if err != nil {
			log.Error().Err(err).Msg("Failed to clear all meetings.")
			return
		}
	}

	if method.ClearAllPreviousMeetings {
		err = ClearAllStudents()
		if err != nil {
			log.Error().Err(err).Msg("Failed to clear all students.")
			return
		}
	}

	failedToAdd := 0
	for _, row := range rows {
		user := User{
			Name: row[0],
			ID:   row[1],
		}

		err = AddNewStudent(user)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to add %s to the DB.", user.Name)
			failedToAdd++
		} else {
			log.Info().Msgf("Added %s to the DB.", user.Name)
		}
	}

	log.Info().Msgf("Added %d students to the DB. Failed to add %d.", len(rows)-failedToAdd, failedToAdd)
}
