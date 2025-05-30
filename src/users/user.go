package users

import (
	"TitanAttendance/src/datastore"
	"TitanAttendance/src/utils"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type User struct {
	ID   string `json:"student_id"`
	Name string `json:"name"`
}

func (u *User) IsValid() error {
	if u.ID == "" {
		return errors.New("student ID is empty")
	}
	if u.Name == "" {
		return errors.New("name is empty")
	}

	for _, c := range u.ID {
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

func (u *User) IDExists() bool {
	for _, v := range GetStudents() {
		if v.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) IsPresent() bool {
	for _, v := range CurrentMeeting.Present {
		if v.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) getFullUserData() bool {
	for _, v := range GetStudents() {
		if v.ID == u.ID {
			u.Name = v.Name
			return true
		}
	}
	return false
}

func (u *User) CheckIn() error {
	if !u.IDExists() {
		return errors.New("student ID does not exist")
	}
	u.getFullUserData()

	conn := datastore.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	meetingCollection := conn.Database(utils.GetDBName()).Collection("meetings")
	findResult := meetingCollection.FindOne(
		ctx,
		map[string]interface{}{
			"date": utils.GetCurrentDate(),
		},
	)
	if errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		CurrentMeeting = Meeting{
			Date:    utils.GetCurrentDate(),
			Absent:  []AbsentStudent{},
			Present: []PresentStudent{},
		}

		for _, student := range GetStudents() {
			if student.ID != u.ID {
				CurrentMeeting.Absent = append(CurrentMeeting.Absent, AbsentStudent{
					ID:   student.ID,
					Name: student.Name,
				})
			}
		}

		_, err := meetingCollection.InsertOne(ctx, CurrentMeeting)
		if err != nil {
			return err
		}
	} else {
		var meeting Meeting
		err := findResult.Decode(&meeting)
		if err != nil {
			return err
		}
		CurrentMeeting = meeting
	}

	if u.IsPresent() {
		return errors.New("already checked in")
	}

	presentStudent := PresentStudent{
		ID:   u.ID,
		Name: u.Name,
		Time: utils.GetCurrentTime(),
	}

	_, err := meetingCollection.UpdateOne(
		ctx,
		map[string]interface{}{
			"date": utils.GetCurrentDate(),
		},
		map[string]interface{}{
			"$pull": map[string]interface{}{
				"absent": AbsentStudent{
					ID:   u.ID,
					Name: u.Name,
				},
			},
			"$push": map[string]interface{}{
				"present": presentStudent,
			},
		},
	)
	if err != nil {
		return err
	}

	for i, v := range CurrentMeeting.Absent {
		if v.ID == u.ID {
			CurrentMeeting.Absent = append(CurrentMeeting.Absent[:i], CurrentMeeting.Absent[i+1:]...)
			break
		}
	}
	CurrentMeeting.Present = append(CurrentMeeting.Present, presentStudent)

	log.Info().Msgf("%s | Checked In!", u.Name)
	return nil
}
