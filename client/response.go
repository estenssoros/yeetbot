package client

import (
	"time"

	"github.com/estenssoros/yeetbot/slack"
	uuid "github.com/satori/go.uuid"
)

// Response stored in elastic
type Response struct {
	ID              uuid.UUID `json:"id"`
	User            string    `json:"user"`
	Report          string    `json:"report"`
	Date            time.Time `json:"date"`
	PendingResponse bool      `json:"pending_response"`
	Responses       []string  `json:"responses"`
}

// GetResponsesByUser gets user response from elastic search
func (c *Client) GetUserResponse(user *slack.User) (*Response, error) {
	return &Response{}, nil
}

// RecordResponse adds response to responses and returns total number of responses recorded
func (c *Client) RecordResponse(user *slack.User, message string) (int, error) {
	return 0, nil
}

// CompleteResponse removes pending status from response and sends "thank you" message to user
func (c *Client) CompleteResponse(user *slack.User) error {
	return nil
}
