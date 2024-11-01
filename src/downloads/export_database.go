package downloads

import (
	"TitanAttendance/src/auth"
	"TitanAttendance/src/users"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
	"net/http"
)

func ExportDatabase(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("TitanAttendancePin")
	if err != nil {
		return
	}

	userAccess, err := auth.CheckWithCookie(*cookie)
	if err != nil {
		return
	}

	if !userAccess.IsAdmin() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	xlxsFile := excelize.NewFile()
	defer func() {
		err = xlxsFile.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}
	}()

	meetings, err := users.GetAllMeetings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	for _, meeting := range meetings {
		_, err = xlxsFile.NewSheet(meeting.Date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}

		err = xlxsFile.SetColWidth(meeting.Date, "B", "B", 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}
		err = xlxsFile.SetCellValue(
			meeting.Date,
			"B2",
			fmt.Sprintf("Present (%d)", len(meeting.Present)),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}

		err = xlxsFile.SetColWidth(meeting.Date, "C", "C", 15)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}
		err = xlxsFile.SetCellValue(
			meeting.Date,
			"C2",
			"Check-in Time",
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}

		err = xlxsFile.SetColWidth(meeting.Date, "E", "E", 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}
		err = xlxsFile.SetCellValue(
			meeting.Date,
			"E2",
			fmt.Sprintf("Absent (%d)", len(meeting.Absent)),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Error().Err(err).Msg("error writing response")
			}
			return
		}

		for i, present := range meeting.Present {
			err = xlxsFile.SetCellValue(meeting.Date, fmt.Sprintf("B%d", i+3), present.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Error().Err(err).Msg("error writing response")
				}
				return
			}
			err = xlxsFile.SetCellValue(meeting.Date, fmt.Sprintf("C%d", i+3), present.Time)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Error().Err(err).Msg("error writing response")
				}
				return
			}
		}

		for i, absent := range meeting.Absent {
			err = xlxsFile.SetCellValue(meeting.Date, fmt.Sprintf("E%d", i+3), absent.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Error().Err(err).Msg("error writing response")
				}
				return
			}
		}
	}

	err = xlxsFile.DeleteSheet("Sheet1")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
		return
	}

	xlxsFile.SetActiveSheet(0)

	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Add("Content-Disposition", "attachment;filename=attendance.xlsx")

	w.WriteHeader(http.StatusOK)
	err = xlxsFile.Write(w)
	if err != nil {
		log.Error().Err(err).Msg("error writing response")
	}
}
