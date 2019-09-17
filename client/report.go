package client

import (
	"context"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
)

// SaveReport saves a report to an elastic index
func (c *Client) SaveReport(report *models.Report) error {
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
func (c *Client) GetOrCreateUserReport(meeting *models.Meeting, user *slack.User) (*models.Report, error) {
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("userID", user.ID))
		q = q.Must(elastic.NewTermQuery("meetingID", meeting.ID))
	}
	reports := []*models.Report{}
	if err := es.GetMany(c.reportIndex, q, &reports); err != nil {
		return nil, err
	}
	if len(reports) == 0 {
		return &models.Report{
			MeetingID: meeting.ID,
			UserID:    user.ID,
			CreatedAt: time.Now(),
		}, nil
	}
	return reports[0], nil
}

// SubmitUserReport sends a user report to the meeting channel
func (c *Client) SubmitUserReport(m *models.Meeting, u *slack.User) error {
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("meetingID", m.ID))
		q = q.Must(elastic.NewTermQuery("userID", u.ID))
	}
	reports := []*models.Report{}
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
