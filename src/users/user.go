package users

import (
	"TitanAttendance/src/datastore"
	"TitanAttendance/src/meetings"
	"TitanAttendance/src/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID   string `json:"id"`
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
	for _, v := range meetings.CurrentMeeting.Present {
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

	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := client.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Msg("Failed to rollback transaction.")
		}
	}(tx, ctx)

	date := utils.GetCurrentDate()
	var absentJSON, presentJSON []byte

	err = tx.QueryRow(ctx,
		`SELECT absent_students, present_students FROM meetings WHERE date = $1`,
		date,
	).Scan(&absentJSON, &presentJSON)

	if errors.Is(err, pgx.ErrNoRows) {
		meetings.CurrentMeeting = meetings.Meeting{
			Date:    utils.GetCurrentDate(),
			Absent:  []meetings.AbsentStudent{},
			Present: []meetings.PresentStudent{},
		}

		for _, student := range GetStudents() {
			if student.ID != u.ID {
				meetings.CurrentMeeting.Absent = append(meetings.CurrentMeeting.Absent, meetings.AbsentStudent{
					ID:   student.ID,
					Name: student.Name,
				})
			}
		}
		absentJSON, err = json.Marshal(meetings.CurrentMeeting.Absent)
		if err != nil {
			return err
		}

		presentJSON, err = json.Marshal(meetings.CurrentMeeting.Present)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO meetings (date, absent_students, present_students) VALUES ($1, $2::jsonb, $3::jsonb)`,
			date,
			absentJSON,
			presentJSON,
		)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}

		if meetings.CurrentMeeting.Date != date {
			meetings.CurrentMeeting = meetings.Meeting{
				Date:    date,
				Absent:  []meetings.AbsentStudent{},
				Present: []meetings.PresentStudent{},
			}

			err = json.Unmarshal(absentJSON, &meetings.CurrentMeeting.Absent)
			if err != nil {
				return err
			}

			err = json.Unmarshal(presentJSON, &meetings.CurrentMeeting.Present)
			if err != nil {
				return err
			}
		}
	}

	if u.IsPresent() {
		return errors.New("already checked in")
	}

	presentStudent := meetings.PresentStudent{
		ID:   u.ID,
		Name: u.Name,
		Time: utils.GetCurrentTime(),
	}
	meetings.CurrentMeeting.Present = append(meetings.CurrentMeeting.Present, presentStudent)

	for i, v := range meetings.CurrentMeeting.Absent {
		if v.ID == u.ID {
			meetings.CurrentMeeting.Absent = append(meetings.CurrentMeeting.Absent[:i], meetings.CurrentMeeting.Absent[i+1:]...)
			break
		}
	}

	absentJSON, err = json.Marshal(meetings.CurrentMeeting.Absent)
	if err != nil {
		return err
	}

	presentJSON, err = json.Marshal(meetings.CurrentMeeting.Present)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE meetings SET absent_students = $2::jsonb, present_students = $3::jsonb WHERE date = $1`,
		date,
		absentJSON,
		presentJSON,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	log.Info().Msgf("%s | Checked In!", u.Name)
	return nil
}
