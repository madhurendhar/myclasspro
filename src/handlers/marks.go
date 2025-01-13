package handlers

import (
	"goscraper/src/helpers"
	"goscraper/src/types"
	"log"
)

func GetMarks(token string) (*types.MarksResponse, error) {
	scraper := helpers.NewAcademicsFetch(token)
	marks, err := scraper.GetMarks()
	if err != nil {
		log.Fatal(err)
	}

	return marks, nil

}
