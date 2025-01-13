package handlers

import (
	"goscraper/src/helpers"
	"goscraper/src/types"
)

func GetCourses(token string) (*types.CourseResponse, error) {
	scraper := helpers.NewCoursePage(token)
	course, err := scraper.GetCourses()
	if err != nil {

	}

	return course, nil
}
