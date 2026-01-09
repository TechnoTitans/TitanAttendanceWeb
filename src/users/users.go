package users

import (
	"TitanAttendance/src/datastore"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Exec(ctx, `INSERT INTO students (id, name) VALUES ($1, $2)`, user.ID, user.Name)
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

	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := client.Query(ctx, `SELECT id, name FROM students`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan student row.")
			continue
		}

		users = append(users, user)
	}

	return users
}

func ClearAllStudents() error {
	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Exec(ctx, `DELETE FROM students`)
	if err != nil {
		return err
	}

	users = []User{}
	log.Info().Msg("Cleared all students.")
	return nil
}
