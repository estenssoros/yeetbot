package views

import (
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func EventHandler(c echo.Context) error {
	req := &slack.EventRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// slack requires that a URL be validated with a "challenge"
	if req.Challenge != "" {
		return c.JSON(http.StatusOK, req.Challenge)
	}

	// what report is this from?
	cli, err := c.(*client.Context).Config.NewClientFromChannel(req.Event.Channel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// what user sent this report?
	user, err := cli.GetUserByID(req.Event.User)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// did user ask for report shortcut
	if req.Event.Text == "report" && !cli.HasUserStartedReport(user) {
		if err := cli.InitiateReport(user); err != nil {
			return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client initiate report"))
		}
		return c.NoContent(http.StatusOK)
	}

	// what stage is this user at?
	question, err := cli.GetRecentQuestion(user)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client get last question"))
	}

	// record a response for a step OR overwrite a response from a previous stage
	err = cli.RecordResponse(&client.RecordResponseInput{
		Question: question,
		User:     user,
		Text:     req.Event.Text,
	})

	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client record response"))
	}

	// report is complete! let's do this!
	if cli.IsReportComplete(user) {
		return errors.Wrap(cli.CompleteReport(user), "client complete report")
	}

	// initiate next stage
	if err := cli.PostNextQuestion(user); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "client post next question"))
	}
	return c.NoContent(http.StatusOK)
}
