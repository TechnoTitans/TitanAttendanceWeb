package users

import (
	"TitanAttendance/src/datastore"
	"TitanAttendance/src/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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

	client := datastore.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := client.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	date := utils.GetCurrentDate()
	var absentJSON, presentJSON []byte

	err = tx.QueryRow(ctx,
		`SELECT absent, present FROM meetings WHERE date = $1`,
		date,
	).Scan(&absentJSON, &presentJSON)

	if errors.Is(err, pgx.ErrNoRows) {
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
		absentJSON, err = json.Marshal(CurrentMeeting.Absent)
		if err != nil {
			return err
		}

		presentJSON, err = json.Marshal(CurrentMeeting.Present)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `INSERT INTO meetings (date, absent, present) VALUES ($1, $2::jsonb, $3::jsonb)`,
			date, absentJSON, presentJSON,
		)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	if u.IsPresent() {
		return errors.New("already checked in")
	}

	presentStudent := PresentStudent{
		ID:   u.ID,
		Name: u.Name,
		Time: utils.GetCurrentTime(),
	}

	for i, v := range CurrentMeeting.Absent {
		if v.ID == u.ID {
			CurrentMeeting.Absent = append(CurrentMeeting.Absent[:i], CurrentMeeting.Absent[i+1:]...)
			break
		}
	}
	CurrentMeeting.Present = append(CurrentMeeting.Present, presentStudent)

	absentJSON, err = json.Marshal(CurrentMeeting.Absent)
	if err != nil {
		return err
	}

	presentJSON, err = json.Marshal(CurrentMeeting.Present)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE meetings SET absent = $2::jsonb, present = $3::jsonb WHERE date = $1`,
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
