package models

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Event records what happened
type Event struct {
	Question string `json:"question"`
	Response string `json:"response"`
	TS       string `json:"eventTS"`
	Color    string `json:"color"`
}

// Report a yeetbot report
type Report struct {
	ID             uuid.UUID `json:"id"`
	MeetingID      string    `json:"meetingID"`
	UserID         string    `json:"userID"`
	ScheduledStart time.Time `json:"scheduledStart"`
	CreatedAt      time.Time `json:"createdAt"`
	Events         []*Event  `json:"events"`
	Done           bool      `json:"done"`
}

// EsType for elasticsvc
func (r Report) EsType() string {
	return `report`
}

func (r Report) String() string {
	ju, _ := json.MarshalIndent(r, "", " ")
	return string(ju)
}

// AddEvent adds an event to a report
func (r *Report) AddEvent(event *Event) {
	for _, e := range r.Events {
		if e.Question == event.Question {
			e.Response = event.Response
			e.TS = event.TS
			return
		}
	}
	r.Events = append(r.Events, event)
}
