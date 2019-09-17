package client

import (
	"context"

	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
)

// CreateMeeting creates a meeting in elastic search
func (c *Client) CreateMeeting(m *models.Meeting) error {
	// get full slack user info for user

	for i, u := range m.Users {
		slackUser, err := c.GetUserByID(u.ID)
		if err != nil {
			return errors.Wrap(err, "get user by id")
		}
		m.Users[i] = slackUser
	}

	{
		m.ScheduledStart = m.Schedule.NextReportDate()
		m.Team = c.Config.Team
		m.SetID()
	}

	es := elasticsvc.New(context.Background())
	if err := es.PutOne(c.meetingIndex, m); err != nil {
		return errors.Wrap(err, "put one")
	}

	// create first report instance for each user
	reports := m.CreateReports()
	if err := es.PutMany(c.reportIndex, &reports); err != nil {
		return errors.Wrap(err, "put many reports")
	}
	return nil
}

// StartMeeting starts a meeting
func StartMeeting(m *models.Meeting) error {
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
func (c *Client) GetActiveMeetingsForUser(user *slack.User) ([]*models.Meeting, error) {
	meetings := []*models.Meeting{}
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("started", true))
		q = q.Must(elastic.NewTermQuery("ended", false))
	}
	if err := es.GetMany(c.meetingIndex, q, &meetings); err != nil {
		return nil, errors.Wrap(err, "es get all")
	}
	userMeetings := []*models.Meeting{}
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

// ListActiveMeetings list the active meetings from a team's config
func (c *Client) ListActiveMeetings() ([]*models.Meeting, error) {
	meetings := []*models.Meeting{}
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("started", true))
		q = q.Must(elastic.NewTermQuery("ended", false))
	}
	if err := es.GetMany(c.meetingIndex, q, &meetings); err != nil {
		return nil, errors.Wrap(err, "es get all")
	}
	return meetings, nil
}

// ListPendingMeetings from a teams config
func (c *Client) ListPendingMeetings() ([]*models.Meeting, error) {
	meetings := []*models.Meeting{}
	es := elasticsvc.New(context.Background())
	q := elastic.NewBoolQuery()
	{
		q = q.Must(elastic.NewTermQuery("started", false))
	}
	if err := es.GetMany(c.meetingIndex, q, &meetings); err != nil {
		return nil, errors.Wrap(err, "es get all")
	}
	return meetings, nil
}

// ListAllMeetings list all meetings from a teams config
func (c *Client) ListAllMeetings() ([]*models.Meeting, error) {
	meetings := []*models.Meeting{}
	es := elasticsvc.New(context.Background())
	if _, err := es.GetAll(c.meetingIndex, &meetings); err != nil {
		return nil, errors.Wrap(err, "es get all")
	}
	return meetings, nil
}
