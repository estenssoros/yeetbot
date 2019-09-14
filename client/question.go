package client

import "github.com/estenssoros/yeetbot/slack"

// Question a question to ask a user
// Color is the attachment color
// If options are given, the question will have a drop down
type Question struct {
	Text    string   `yaml:"text"`
	Color   string   `yaml:"color"`
	Options []string `yaml:"options,omitempty"`
}

// PostFirstQuestion sends the first question to the user
// and creates a new response in elastic search with "pending_response"
func (c *Client) PostFirstQuestion(user *slack.User, response *Response) error {
	return nil
}

// PostNextQuestion sends the next question to the user
// and sets status to pending response
func (c *Client) PostNextQuestion(user *slack.User, response *Response) error {
	return nil
}
