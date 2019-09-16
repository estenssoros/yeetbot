package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
)

// Meeting holds all the info for a given day's meeting
type Meeting struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Channel        string      `json:"channel"`
	Team           string      `json:"team"`
	ScheduledStart time.Time   `json:"scheduledStart"`
	Users          []*User     `json:"users"`
	Questions      []*Question `json:"questions"`
	Started        bool        `json:"started"`
	Ended          bool        `json:"ended"`
	IntroMessage   string      `json:"introMessage"`
	Schedule       *Schedule   `json:"-"`
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

// CreateMeeting creates a meeting in elastic search
func (c *Client) CreateMeeting(m *Meeting) error {
	es := elasticsvc.New(context.Background())
	{
		m.ScheduledStart = m.Schedule.NextReportDate()
		m.Team = c.Config.Team
		m.SetID()
	}
	if err := es.PutOne("yeetmeet", m); err != nil {
		return errors.Wrap(err, "put one")
	}
	return nil
}

// StartMeeting starts a meeting
func StartMeeting(m *Meeting) error {
	config, err := ConfigFromEnv()
	if err != nil {
		return errors.Wrap(err, "config from env")
	}

	client := config.NewClient()

	for _, u := range m.Users {
		if err := client.SendGreeting(m, u); err != nil {
			return errors.Wrap(err, "client send greeting")
		}
		if err := client.PostFirstQuestion(m, u); err != nil {
			return errors.Wrap(err, "client post first question")
		}
	}
	return nil
}

// GetActiveMeetingsForUser gets the active meetings for a user
func (c *Client) GetActiveMeetingsForUser(user *slack.User) ([]*Meeting, error) {
	meetings := []*Meeting{}
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("started", true))
		q = q.Must(elastic.NewTermQuery("ended", false))
	}
	if err := es.GetMany("yeetmeet", q, &meetings); err != nil {
		return nil, errors.Wrap(err, "es get all")
	}
	userMeetings := []*Meeting{}
	for _, meeting := range meetings {
		if meeting.HasUser(user) {
			userMeetings = append(userMeetings, meeting)
		}
	}
	if len(userMeetings) > 0 {
		return userMeetings, nil
	}
	return nil, errors.New("user has no active meetings")
}
