package views

import (
	"context"
	"net/http"
	"time"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
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
	for i, m := range messages {
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

	logrus.Infof("event type: %s", req.Type)

	// what report is this from?
	cli, err := c.(*client.Context).NewClientFromUser(req.Event.User)
	if err != nil {
		logrus.Error(errors.Wrap(err, "new client from channel"))
		return c.JSON(http.StatusInternalServerError, err)
	}

	es := elasticsvc.New(context.Background())

	messages, err := cli.ListTodayMessages(req.Event.Channel)
	if err != nil {
		logrus.Error(errors.Wrap(err, "cli list messages"))
		return c.JSON(http.StatusInternalServerError, err)
	}

	msgIdx, err := getMessageIndexByTS(messages, req.Event.Ts)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	question, err := getMessageByUserID(message[msgIdx:], cli.YeetUser)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	doc := &client.Response{
		Team:     req.TeamID,
		UserID:   req.Event.User,
		EventTS:  req.Event.EventTs,
		Date:     time.Now(),
		Text:     req.Event.Text,
		Question: question.Text,
	}

	if err := es.PutOne("yeetbot", doc); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "es put one"))
	}
	logrus.Infof("added %s", doc.ID)

	// // did user ask for report shortcut
	// if req.Event.Text == "report" && !cli.HasUserStartedReport(user) {
	// 	if err := cli.InitiateReport(user); err != nil {
	// 		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client initiate report"))
	// 	}
	// 	return c.NoContent(http.StatusOK)
	// }

	// // what stage is this user at?
	// question, err := cli.GetRecentQuestion(user)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client get last question"))
	// }

	// // record a response for a step OR overwrite a response from a previous stage
	// err = cli.RecordResponse(&client.RecordResponseInput{
	// 	Question: question,
	// 	User:     user,
	// 	Text:     req.Event.Text,
	// })

	// if err != nil {
	// 	logrus.Error(err)
	// 	return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client record response"))
	// }

	// // report is complete! let's do this!
	// if cli.IsReportComplete(user) {
	// 	return errors.Wrap(cli.CompleteReport(user), "client complete report")
	// }

	// // initiate next stage
	// if err := cli.PostNextQuestion(user); err != nil {
	// 	logrus.Error(err)
	// 	return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client post next question"))
	// }
	return c.NoContent(http.StatusOK)
}
