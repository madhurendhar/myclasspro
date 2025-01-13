package handlers

import (
	"goscraper/src/helpers"
	"goscraper/src/types"
	"log"
	"time"
)

func GetCalendar(token string) (*types.CalendarResponse, error) {
	scraper := helpers.NewCalendarFetcher(time.Now(), token)
	calendar, err := scraper.GetCalendar()
	if err != nil {
		log.Fatal(err)
	}

	return calendar, nil

}
