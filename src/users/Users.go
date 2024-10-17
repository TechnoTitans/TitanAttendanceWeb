package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"os"
)

var users []User

func AddNewStudent(user User) error {
	err := user.IsValid()
	if err != nil {
		return err
	}

	collection := database.GetFireDB().Client.Collection("TitanAttendance")
	_, err = collection.Doc("students").Set(
		context.Background(),
		map[string]interface{}{
			user.StudentID: user.Name,
		},
		firestore.MergeAll,
	)
	if err == nil {
		users = append(users, user)
	}

	return err
}

func GetStudents() []User {
	if len(users) > 0 {
		return users
	}

	log.Info().Msg("Getting students from the database.")

	collection := database.GetFireDB().Client.Collection("TitanAttendance")
	doc, err := collection.Doc("students").Get(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get students.")
		return nil
	}

	data := doc.Data()
	for k, v := range data {
		users = append(users, User{
			StudentID: k,
			Name:      v.(string),
		})
	}

	return users
}

func ClearAllStudents() error {
	collection := database.GetFireDB().Client.Collection("TitanAttendance")
	_, err := collection.Doc("students").Delete(context.Background())
	if err != nil {
		return err
	}

	users = []User{}
	return nil
}

func ClearAllMeetings() error {
	collection := database.GetFireDB().Client.Collection("TitanAttendance")
	_, err := collection.Doc("meetings").Delete(context.Background())
	if err != nil {
		return err
	}
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
			Name:      row[0],
			StudentID: row[1],
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
