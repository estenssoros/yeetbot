package views

import (
	"log"
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
)

func EventHandler(cc client.Context) error {
	req := &slack.EventRequest{}
	if err := cc.Bind(req); err != nil {
		return cc.JSON(http.StatusBadRequest, err)
	}
	if req.Challenge != "" {
		return cc.JSON(http.StatusOK, req.Challenge)
	}
	c := cc.Config.NewClient(&client.Report{})
	user, err := c.GetUserFromRequest(req)
	if err != nil {
		log.Println(err)
		return cc.JSON(http.StatusInternalServerError, err)
	}
	report, err := client.FindReportByUser(user, cc.UserReports)
	if err != nil {
		log.Println(err)
		return cc.JSON(http.StatusInternalServerError, err)
	}
	c.Report = report
	response, err := c.GetUserResponse(user)
	if err != nil {
		log.Println(err)
		return cc.JSON(http.StatusInternalServerError, err)
	}
	if response == nil {
		err := c.PostFirstQuestion(user)
		if err != nil {
			log.Println(err)
			return cc.JSON(http.StatusInternalServerError, err)
		}
		cc.JSON(http.StatusOK, req.Challenge)
	}
	c.Response = response
	total, err := c.RecordResponse(user, req.Event.Text)
	if err != nil {
		log.Println(err)
		return cc.JSON(http.StatusInternalServerError, err)
	}
	if total == len(c.Report.Questions) {
		err := c.PostNextQuestion(user)
		if err != nil {
			log.Println(err)
			return cc.JSON(http.StatusInternalServerError, err)
		}
		cc.JSON(http.StatusOK, req.Challenge)
	}
	c.CompleteResponse(user)
	// send complete message to report channel

	return cc.JSON(http.StatusOK, req.Challenge)
}
