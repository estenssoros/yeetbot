package views

import (
	"fmt"
	"log"
	"net/http"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/models"
	"github.com/labstack/echo"
)

func EventHandler(c echo.Context) error {
	req := &models.EventRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if req.Challenge != "" {
		return c.JSON(http.StatusOK, req.Challenge)
	}
	client, err := client.NewAWS()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	messages, err := client.ListMessages(req.Event.Channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(messages)
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
