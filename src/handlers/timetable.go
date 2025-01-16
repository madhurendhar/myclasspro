package handlers

import (
	"goscraper/src/helpers"
	"goscraper/src/types"
)

func GetTimetable(token string) (*types.TimetableResult, error) {
	scraper := helpers.NewTimetable(token)
	timetable, err := scraper.GetTimetable()

	return timetable, err
}
