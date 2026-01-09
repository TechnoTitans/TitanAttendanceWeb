package users

import (
	"TitanAttendance/src/datastore"
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
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
	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Exec(ctx, `DELETE FROM meetings`)
	if err != nil {
		return err
	}

	log.Info().Msg("Cleared all meetings.")
	return nil
}

func GetAllMeetings() ([]Meeting, error) {
	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := client.Query(ctx, `SELECT date, absent, present FROM meetings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []Meeting
	for rows.Next() {
		var meeting Meeting
		var absentJSON, presentJSON []byte

		err := rows.Scan(&meeting.Date, &absentJSON, &presentJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(absentJSON, &meeting.Absent); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(presentJSON, &meeting.Present); err != nil {
			return nil, err
		}

		meetings = append(meetings, meeting)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return meetings, nil
}
