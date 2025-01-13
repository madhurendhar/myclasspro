package types

type Slot struct {
	Day      int      `json:"day"`
	DayOrder string   `json:"dayOrder"`
	Slots    []string `json:"slots"`
}

type Batch struct {
	Batch string `json:"batch"`
	Slots []Slot `json:"slots"`
}

type DaySchedule struct {
	Day   int      `json:"day"`
	Table []string `json:"table"`
}

type TimetableResult struct {
	RegNumber string        `json:"regNumber"`
	Batch     string        `json:"batch"`
	Schedule  []DaySchedule `json:"schedule"`
}