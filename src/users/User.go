package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
}

var absentExists = false

func (u *User) IsValid() error {
	if u.StudentID == "" {
		return errors.New("student ID is empty")
	}
	if u.Name == "" {
		return errors.New("name is empty")
	}

	for _, c := range u.StudentID {
		if c < '0' || c > '9' {
			return errors.New("student ID contains non-numeric characters")
		}
	}

	for _, c := range u.Name {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && c != ' ' && c != '-' {
			return errors.New("name contains invalid characters. name must contain a letter space or hyphen")
		}
	}

	return nil
}

func (u *User) Exists() bool {
	for _, v := range GetStudents() {
		if v.StudentID == u.StudentID {
			u.Name = v.Name
			return true
		}
	}
	return false
}

func (u *User) CheckIn() error {
	if !u.Exists() {
		return errors.New("student ID does not exist")
	}

	collection := database.GetFireDB().Client.Collection("TitanAttendance")
	currentMeeting := collection.Doc("meetings").Collection(utils.GetCurrentDate())

	_, err := currentMeeting.Doc("Present").Set(
		context.Background(),
		map[string]interface{}{
			u.Name: utils.GetCurrentTime(),
		},
		firestore.MergeAll,
	)
	if err != nil {
		return errors.New("failed to add student to present list in Firebase")
	}

	if !absentExists {
		absentExists = true
		docs, err := currentMeeting.Doc("Absent").Get(context.Background())
		if err != nil && status.Code(err) != codes.NotFound {
			return errors.New("failed to get documents from this meeting in Firebase")
		}
		if !docs.Exists() {
			for _, student := range GetStudents() {
				if student.StudentID != u.StudentID {
					_, err = currentMeeting.Doc("Absent").Set(
						context.Background(),
						map[string]interface{}{
							student.Name: "",
						},
						firestore.MergeAll,
					)
					if err != nil {
						return errors.New("failed to create absent list in Firebase")
					}
				}
			}
		}
	}

	log.Info().Msgf("%s | Checked In!", u.Name)
	return nil
}
