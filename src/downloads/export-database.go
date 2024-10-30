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
		if err = xlxsFile.Close(); err != nil {
			log.Error().Err(err).Msg("error closing file")
		}
	}()

	meetings, err := users.GetAllMeetings()
	if err != nil {
		log.Error().Err(err).Msg("error getting all meetings")
		return
	}

	for _, date := range meetings {
		fmt.Println(date.Date)
		_, err = xlxsFile.NewSheet("date")
		if err != nil {
			log.Error().Err(err).Msg("error creating sheet")
		}

	}

	index, err := xlxsFile.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set value of a cell.
	err = xlxsFile.SetCellValue("Sheet2", "A2", "Hello world.")
	if err != nil {
		log.Error().Err(err).Msg("error setting cell value")
	}
	err = xlxsFile.SetCellValue("Sheet1", "B2", 100)
	if err != nil {
		log.Error().Err(err).Msg("error setting cell value")
	}

	// Set active sheet of the workbook.
	xlxsFile.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err = xlxsFile.SaveAs("Book1.xlsx"); err != nil {
		log.Error().Err(err).Msg("error saving file")
	}
}
