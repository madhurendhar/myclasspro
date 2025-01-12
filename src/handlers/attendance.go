package handlers

import (
	"goscraper/src/helpers"
	"log"
	"goscraper/src/types"
)

func GetAttendance(token string) (*types.AttendanceResponse, error) {
	scraper := helpers.NewAcademicsFetch(token)
	attendance, err := scraper.GetAttendance()
	if err != nil {
		log.Fatal(err)
	}

	return attendance, nil


}
