package views

import (
	"log"
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
)

func EventHandler(c client.Context) error {
	req := &slack.EventRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// slack requires that a URL be validated with a "challenge"
	if req.Challenge != "" {
		return c.JSON(http.StatusOK, req.Challenge)
	}

	// what report is this from (maybe from report channel?)
	// what user sent this report?
	// what stage is this user at?
	// record response
	// initiate next stage

	client := c.Config.NewClient(&client.Report{})
	user, err := client.GetUserFromRequest(req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	report, err := client.FindReportByUser(user, c.UserReports)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	client.Report = report
	response, err := client.GetUserResponse(user)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if response == nil {
		if err := client.PostFirstQuestion(user); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, req.Challenge)
	}
	client.Response = response
	total, err := client.RecordResponse(user, req.Event.Text)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if total == len(client.Report.Questions) {
		if err := client.PostNextQuestion(user); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, req.Challenge)
	}
	client.CompleteResponse(user)
	// send complete message to report channel

	return c.NoContent(http.StatusOK)
}
