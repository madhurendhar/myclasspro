package types

type Day struct {
	Date     string `json:"date"`
	Day      string `json:"day"`
	Event    string `json:"event"`
	DayOrder string `json:"dayOrder"`
}

type CalendarMonth struct {
	Month string `json:"month"`
	Days  []Day  `json:"days"`
}

type CalendarResponse struct {
	Today    *Day            `json:"today"`
	Index    int             `json:"index"`
	Calendar []CalendarMonth `json:"calendar"`
	Status   int             `json:"status"`
	Error    bool            `json:"error,omitempty"`
	Message  string          `json:"message,omitempty"`
}
