package views

import (
	"fmt"
	"log"
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/labstack/echo"
)

func EventHandler(c echo.Context) error {
	req := &slack.EventRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if req.Challenge != "" {
		return c.JSON(http.StatusOK, req.Challenge)
	}
	cc := c.(*client.Context)
	client := cc.Config.NewClient(&client.Report{})
	user, err := client.GetUserFromRequest(req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	report, err := cc.FindUserReport(user.RealName)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	client.Report = report
	fmt.Println(client.Report)
	// TODO: find messages since report, what message step are we on record response

	// client, err := client.NewFromReader()
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// messages, err := client.ListMessages(req.Event.Channel)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// fmt.Println(messages)
	// user, err := client.GetUserFromRequest(req)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// lastMsg, err := client.GetLastMessageFromUser(user)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// fmt.Println(lastMsg)
	// if err := client.UpdateUserFlow(user, lastMsg, req); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	return c.JSON(http.StatusOK, req.Challenge)
}
