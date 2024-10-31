package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type User struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
}

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

	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findResult := conn.Database(utils.DBName).Collection("meetings").FindOne(
		ctx,
		map[string]interface{}{
			"date": utils.GetCurrentDate(),
		},
	)
	if errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		meeting := Meeting{
			Date:    utils.GetCurrentDate(),
			Absent:  []AbsentStudent{},
			Present: []PresentStudent{},
		}

		for _, student := range GetStudents() {
			if student.StudentID != u.StudentID {
				meeting.Absent = append(meeting.Absent, AbsentStudent{StudentName: student.Name})
			}
		}

		_, err := conn.Database(utils.DBName).Collection("meetings").InsertOne(ctx, meeting)
		if err != nil {
			return errors.New("failed to create new meeting")
		}
	}

	_, err := conn.Database(utils.DBName).Collection("meetings").UpdateOne(
		ctx,
		map[string]interface{}{
			"date": utils.GetCurrentDate(),
		},
		map[string]interface{}{
			"$pull": map[string]interface{}{
				"absent": AbsentStudent{StudentName: u.Name},
			},
			"$push": map[string]interface{}{
				"present": PresentStudent{
					StudentName: u.Name,
					Time:        utils.GetCurrentTime(),
				},
			},
		},
	)
	if err != nil {
		return errors.New("failed to check in")
	}

	log.Info().Msgf("%s | Checked In!", u.Name)
	return nil
}
