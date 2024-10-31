package utils

import (
	"fmt"
	"strings"
)

type CSVUploadMethod struct {
	ClearAllPreviousIDs      bool
	ClearAllPreviousMeetings bool
}

func AskForNewCSVMethod() (CSVUploadMethod, error) {
	var uploadMethod = CSVUploadMethod{}

	fmt.Println("Would you like to erase all previous student IDs and replace them with the ones in the csv? " +
		"(y/n)")
	for {
		clearAllPreviousIDs, err := GetUserInput()
		if err != nil {
			return CSVUploadMethod{}, err
		}
		clearAllPreviousIDs = strings.ToLower(clearAllPreviousIDs)
		if clearAllPreviousIDs != "y" && clearAllPreviousIDs != "n" {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		} else {
			uploadMethod.ClearAllPreviousMeetings = clearAllPreviousIDs == "y"
			break
		}
	}

	fmt.Println("Would you like to clear all previous meeting attendance data? (y/n)")
	for {
		clearAllPreviousMeetings, err := GetUserInput()
		if err != nil {
			return CSVUploadMethod{}, err
		}
		clearAllPreviousMeetings = strings.ToLower(clearAllPreviousMeetings)
		if clearAllPreviousMeetings != "y" && clearAllPreviousMeetings != "n" {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		} else {
			uploadMethod.ClearAllPreviousMeetings = clearAllPreviousMeetings == "y"
			break
		}
	}

	return uploadMethod, nil
}
