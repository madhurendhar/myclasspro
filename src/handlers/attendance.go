package handlers

import (
	"goscraper/src/helpers"
	"goscraper/src/types"
	"log"
)

func GetAttendance(token string) (*types.AttendanceResponse, error) {
	scraper := helpers.NewAcademicsFetch(token)
	attendance, err := scraper.GetAttendance()
	if err != nil {
		log.Fatal(err)
	}

	return attendance, nil

}
