package views

import (
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getMessageIndexByTS(messages []*slack.HistoryMessage, eventTS string) (int, error) {
	for i, m := range messages {
		if m.Ts == eventTS {
			return i, nil
		}
	}
	return 0, errors.New("message not found")
}

func getMessageByUserID(messages []*slack.HistoryMessage, userID string) (*slack.HistoryMessage, error) {
	for _, m := range messages {
		if m.User == userID {
			return m, nil
		}
	}
	return nil, errors.New("message not found")
}

func EventHandler(c echo.Context) error {
	req := &slack.EventRequest{}
	if err := c.Bind(req); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// slack requires that a URL be validated with a "challenge"
	if req.Challenge != "" {
		return c.JSON(http.StatusOK, req.Challenge)
	}

	cli := c.(*client.Context).NewEmptyClient()

	user, err := cli.GetUserByID(req.Event.User)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	activeMeetings, err := cli.GetActiveMeetingsForUser(user)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	meeting := activeMeetings[0]

	for _, m := range activeMeetings {
		if m.ScheduledStart.After(meeting.ScheduledStart) {
			meeting = m
		}
	}

	lastMessage, err := cli.GetLastYeetQuestion(req.Event.Channel)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	report, err := cli.GetOrCreateUserReport(meeting, user)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	event := &client.Event{
		Question: lastMessage.Text,
		Response: req.Event.Text,
		TS:       req.Event.EventTs,
	}

	report.AddEvent(event)

	if err := cli.SaveReport(report); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := cli.NextStage(meeting, lastMessage, user); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
