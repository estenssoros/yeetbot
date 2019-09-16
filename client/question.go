package client

import (
	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
)

// Question a question to ask a user
// Color is the attachment color
// If options are given, the question will have a drop down
type Question struct {
	Text    string   `yaml:"text"`
	Color   string   `yaml:"color"`
	Options []string `yaml:"options,omitempty"`
}

// PostFirstQuestion sends the first question to the user
func (c *Client) PostFirstQuestion(m *Meeting, u *User) error {
	msg := &slack.Message{
		Text:    m.Questions[0].Text,
		Channel: "@" + u.Name,
		AsUser:  true,
	}
	return errors.Wrap(c.SendMessage(msg), "send message")
}

// NextStage sends the next question to the user
func (c *Client) NextStage(m *Meeting, q *Question, u *slack.User) error {
	var questionIdx int
	for idx, question := range m.Questions {
		if question.Text == q.Text {
			questionIdx = idx
			break
		}
	}
	if questionIdx == len(m.Questions)-1 {

		return errors.Wrap(c.SubmitUserReport(m, u), "submit user report")
	}
	return errors.Wrap(c.PostQuestion(m.Questions[questionIdx+1], u), "post next question")
}

// PostQuestion sends a question to a slack user
func (c *Client) PostQuestion(q *Question, u *slack.User) error {
	msg := &slack.Message{
		Text:    q.Text,
		Channel: "@" + u.Name,
		AsUser:  true,
	}
	return errors.Wrap(c.SendMessage(msg), "send message")
}

// GetLastYeetQuestion gets the last question asked by yeetbot
func (c *Client) GetLastYeetQuestion(channelID string) (*Question, error) {
	messages, err := c.ListMessages(channelID)
	if err != nil {
		return nil, err
	}
	for _, m := range messages {
		if m.User == c.YeetUser {
			return &Question{
				Text: m.Text,
			}, nil
		}
	}
	return nil, errors.New("could not locate last yeet message")
}
