package client

import (
	"time"

	"github.com/estenssoros/yeetbot/slack"
	uuid "github.com/satori/go.uuid"
)

// Response stored in elastic
type Response struct {
	ID       uuid.UUID `json:"id"`
	Team     string    `json:"team"`
	Report   string    `json:"report"`
	UserID   string    `json:"user_id"`
	EventTS  string    `json:"event_ts"`
	Date     time.Time `json:"date"`
	Question string    `json:"question"`
	Text     string    `json:"text"`
}

func (r Response) EsType() string {
	return `response`
}

type RecordResponseInput struct {
	Question *Question
	User     *slack.User
	Text     string
}

// GetUserResponse gets user response from elastic search
// list most recent responses on channel
func (c *Client) GetUserResponse(user *slack.User) (*Response, error) {
	return &Response{}, nil
}

// RecordResponse adds response to responses and returns total number of responses recorded
func (c *Client) RecordResponse(input *RecordResponseInput) error {
	return nil
}

// CompleteResponse removes pending status from response and sends "thank you" message to user
func (c *Client) CompleteResponse(user *slack.User) error {
	return nil
}
