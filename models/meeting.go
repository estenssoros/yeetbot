package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	uuid "github.com/satori/go.uuid"
)

// Meeting holds all the info for a given day's meeting
type Meeting struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Channel        string        `json:"channel"`
	Team           string        `json:"team"`
	ScheduledStart time.Time     `json:"scheduledStart"`
	Users          []*slack.User `json:"users"`
	Questions      []*Question   `json:"questions"`
	Started        bool          `json:"started"`
	Ended          bool          `json:"ended"`
	IntroMessage   string        `json:"introMessage"`
	Schedule       *Schedule     `json:"-"`
}

// SetID crafts a unique id for a meeting
func (m *Meeting) SetID() {
	m.ID = fmt.Sprintf("%s-%s-%s", m.Team, m.Name, m.ScheduledStart.Format("20060102"))
}

func (m Meeting) String() string {
	ju, _ := json.Marshal(m)
	return string(ju)
}

// EsType for elastic search service
func (m Meeting) EsType() string {
	return `meeting`
}

// HasUser loop through report to see if has user
func (m *Meeting) HasUser(u *slack.User) bool {
	for _, user := range m.Users {
		if user.ID == u.ID {
			return true
		}
	}
	return false
}

// CraftEvents creates events from question
func (m *Meeting) CraftEvents() []*Event {
	events := []*Event{}
	for _, q := range m.Questions {
		event := &Event{
			Question: q.Text,
			Color:    q.Color,
		}
		events = append(events, event)
	}
	return events
}

// CreateReports creatres reports from users
func (m *Meeting) CreateReports() []*Report {
	reports := []*Report{}
	for _, user := range m.Users {
		report := &Report{
			ID:             uuid.Must(uuid.NewV4()),
			MeetingID:      m.ID,
			UserID:         user.ID,
			ScheduledStart: m.Schedule.UserReportDate(user),
			CreatedAt:      time.Now(),
			Events:         m.CraftEvents(),
		}
		reports = append(reports, report)
	}
	return reports
}
