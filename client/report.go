package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/seaspancode/services/elasticsvc"
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

// SaveReport saves a report to an elastic index
func (c *Client) SaveReport(report *Report) error {
	es := elasticsvc.New(context.Background())
	if err := es.PutOne(c.reportIndex, report); err != nil {
		return errors.Wrap(err, "put one")
	}
	return nil
}

// InitiateReport initiates a new report for a user
func (c *Client) InitiateReport(user *slack.User) error {
	// TODO this
	return nil
}

// IsReportComplete checks to see if the user has completed all questions
func (c *Client) IsReportComplete(user *slack.User) bool {
	// TODO this
	return true
}

// CompleteReport sends the users report to slack
func (c *Client) CompleteReport(user *slack.User) error {
	// TODO this
	return nil
}

// GetOrCreateUserReport gets a user report if it exists. Creates a new one if it doesn't exist
func (c *Client) GetOrCreateUserReport(meeting *Meeting, user *slack.User) (*Report, error) {
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("userID", user.ID))
		q = q.Must(elastic.NewTermQuery("meetingID", meeting.ID))
	}
	reports := []*Report{}
	if err := es.GetMany(c.reportIndex, q, &reports); err != nil {
		return nil, err
	}
	if len(reports) == 0 {
		return &Report{
			MeetingID: meeting.ID,
			UserID:    user.ID,
			CreatedAt: time.Now(),
		}, nil
	}
	return reports[0], nil
}

// SubmitUserReport sends a user report to the meeting channel
func (c *Client) SubmitUserReport(m *Meeting, u *slack.User) error {
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("meetingID", m.ID))
		q = q.Must(elastic.NewTermQuery("userID", u.ID))
	}
	reports := []*Report{}
	if err := es.GetMany(c.reportIndex, q, &reports); err != nil {
		return errors.Wrap(err, "es get many")
	}

	if want, have := 1, len(reports); want != have {
		return errors.Errorf("fetching reports: wanted %d, have %d", want, have)
	}
	report := reports[0]
	attachments := []*slack.Attachment{}
	for _, event := range report.Events {
		attachments = append(attachments, &slack.Attachment{
			Title: event.Question,
			Text:  event.Response,
		})
	}
	msg := &slack.Message{
		Text:        fmt.Sprintf("*%s* posted an update for *Daily Standup*", u.Name),
		Channel:     m.Channel,
		Attachments: attachments,
	}
	if err := c.SendMessage(msg); err != nil {
		return errors.Wrap(err, "send message")
	}
	report.Done = true
	return errors.Wrap(es.PutOne(c.reportIndex, report), "put report")
}
