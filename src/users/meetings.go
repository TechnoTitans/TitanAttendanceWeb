package users

import (
	"TitanAttendance/src/database"
	"TitanAttendance/src/utils"
	"context"
	"time"
)

type AbsentStudent struct {
	StudentName string `json:"student_name"`
}

type PresentStudent struct {
	StudentName string `json:"student_name"`
	Time        string `json:"time"`
}

type Meeting struct {
	Date    string           `json:"date"`
	Absent  []AbsentStudent  `json:"absent"`
	Present []PresentStudent `json:"present"`
}

func ClearAllMeetings() error {
	conn := database.GetConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Database(utils.DBName).Collection("meetings").DeleteMany(ctx, nil)
	if err != nil {
		return err
	}

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
