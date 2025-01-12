package handlers

import (
	"goscraper/src/helpers"
	"log"
	"goscraper/src/types"
)

func GetMarks(token string) (*types.MarksResponse, error) {
	scraper := helpers.NewAcademicsFetch(token)
	marks, err := scraper.GetMarks()
	if err != nil {
		log.Fatal(err)
	}

	return marks, nil


}
