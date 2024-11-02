package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type AbsentStudent struct {
	ID   string `json:"student_id"`
	Name string `json:"student_name"`
}

type PresentStudent struct {
	ID   string `json:"student_id"`
	Name string `json:"student_name"`
	Time string `json:"time"`
}

type Meeting struct {
	Date    string           `json:"date"`
	Absent  []AbsentStudent  `json:"absent"`
	Present []PresentStudent `json:"present"`
}

var CurrentMeeting Meeting

func ClearAllMeetings() error {
	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Database(utils.DBName).Collection("meetings").DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	log.Info().Msg("Cleared all meetings.")
	return nil
}

func GetAllMeetings() ([]Meeting, error) {
	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	meetingCollection := conn.Database(utils.DBName).Collection("meetings")
	cursor, err := meetingCollection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var meetings []Meeting
	for cursor.Next(ctx) {
		var meeting Meeting
		err = cursor.Decode(&meeting)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}
