package client

import (
	"encoding/json"

	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
)

var (
	defaultUsername = "yeetbot"
	defaultIcon     = ":ghost:"
	yeetbotBucket   = "yeetbot"
	userIDMap       map[string]*string
)

func init() {
	userIDMap = map[string]*string{}
}

// SlackInterface implements all slack operations
type SlackInterface interface {
	SendMessage(*slack.Message) error
	ListUsers() ([]*slack.User, error)
	ListChannels() ([]*slack.Channel, error)
	ListMessages(string) ([]*slack.HistoryMessage, error)
	ListTodayMessages(string) ([]*slack.HistoryMessage, error)
	ListDirectMessageChannels() ([]*slack.Channel, error)
	DeleteBotMessage(string, string) error
	GetUserByID(string) (*slack.User, error)
	SetVerbose(bool)
}

// ElasticInterface implements all elastic operations
type ElasticInterface interface {
	PutOne(string, interface{}) error
	GetMany(string, elastic.Query, interface{}) error
	PutMany(string, interface{}) error
	GetAll(string, interface{}) (*elasticsvc.Result, error)
}

// Client the guy that does all the work
type Client struct {
	UserToken    string `json:"user_token"`
	BotToken     string `json:"bot_token"`
	YeetUser     string `json:"yeet_user"`
	Debug        bool   `json:"debug"`
	meetingIndex string
	reportIndex  string
	*Config
	slack   SlackInterface
	elastic ElasticInterface
}

func (c Client) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// SendMessage sends a slack message
func (c *Client) SendMessage(msg *slack.Message) error {
	return errors.Wrap(c.slack.SendMessage(msg), "slack send message")
}

// GenericMessage sends a generic message
func (c *Client) GenericMessage(u *slack.User, text string) error {
	msg := &slack.Message{
		Text:    text,
		Channel: "@" + u.Name,
		AsUser:  true,
	}
	return errors.Wrap(c.SendMessage(msg), "client send message")
}

// SendGreeting crafts and sends the greeting message
func (c *Client) SendGreeting(m *models.Meeting, user *slack.User) error {
	if c.BotToken == "" {
		return errors.New("missing bot token")
	}
	if c.UserToken == "" {
		return errors.New("missing user token")
	}
	text, err := user.Template(m.IntroMessage)
	if err != nil {
		return errors.Wrap(err, "user template")
	}
	msg := &slack.Message{
		Text:    text,
		Channel: "@" + user.Name,
		AsUser:  true,
		Attachments: []*slack.Attachment{
			&slack.Attachment{
				Text:       "When you are ready, please answer the following question:",
				Color:      "#3AA3E3",
				CallbackID: "quick-reply",
				Actions: []*slack.Action{
					&slack.Action{
						Name: "quick-reply",
						Text: "quick reply",
						Type: "select",
						Options: []*slack.Option{
							&slack.Option{
								Text:  "Skip",
								Value: "skip",
							},
							&slack.Option{
								Text:  "Same as last time",
								Value: "same",
							},
							&slack.Option{
								Text:  "I'm on vacation 🏝",
								Value: "vacation",
							},
						},
					},
				},
				AttachmentType: "default",
			},
		},
	}
	if err := c.SendMessage(msg); err != nil {
		return errors.Wrap(err, "send message")
	}
	return nil
}

// ListUsers list users in workspace
func (c *Client) ListUsers() ([]*slack.User, error) {
	return c.slack.ListUsers()
}

// ListChannels list channels in workspace
func (c *Client) ListChannels() ([]*slack.Channel, error) {
	return c.slack.ListChannels()
}

// ListMessages lists messages in channel
func (c *Client) ListMessages(channelID string) ([]*slack.HistoryMessage, error) {
	return c.slack.ListMessages(channelID)
}

// ListTodayMessages lists messages from today
func (c *Client) ListTodayMessages(channelID string) ([]*slack.HistoryMessage, error) {
	return c.slack.ListTodayMessages(channelID)
}

// ListDirectMessageChannels lists direct message channels
func (c *Client) ListDirectMessageChannels() ([]*slack.Channel, error) {
	return c.slack.ListDirectMessageChannels()
}

// DeleteBotMessage deletes a bot message
func (c *Client) DeleteBotMessage(channelID string, messageTS string) error {
	return c.slack.DeleteBotMessage(channelID, messageTS)
}

// GetUserByName lists slack users and then returns the user
func (c *Client) GetUserByName(userName string) (*slack.User, error) {
	users, err := c.ListUsers()
	if err != nil {
		return nil, errors.Wrap(err, "client list users")
	}
	for _, u := range users {
		if u.Name == userName || u.RealName == userName {
			userIDMap[userName] = &u.ID
			return u, nil
		}
	}
	return nil, errors.New("could not locate user")
}

// GetUserByID searches for a user by iD
func (c *Client) GetUserByID(userID string) (*slack.User, error) {
	return c.slack.GetUserByID(userID)
}
