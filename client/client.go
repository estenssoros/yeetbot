package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/estenssoros/yeetbot/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	slackurl        = "https://slack.com/api"
	defaultUsername = "yeetbot"
	defaultIcon     = ":ghost:"
	yeetbotBucket   = "yeetbot"
	userIDMap       map[string]*string
)

func init() {
	userIDMap = map[string]*string{}
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
}

func (c Client) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// SendMessage sends a slack message
func (c *Client) SendMessage(msg *slack.Message) error {
	u := slackurl + "/" + slack.ChatPostMessage

	if err := c.postRequest(u, msg); err != nil {
		return errors.Wrap(err, "client post request")
	}
	return nil
}

func (c *Client) postRequest(url string, v interface{}) error {

	logrus.Infof("%s %s", http.MethodPost, url)

	ju, _ := json.Marshal(v)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(ju))

	req.Header.Set("Authorization", "Bearer "+c.BotToken)
	req.Header.Set("Content-Type", "application/json;charset=iso-8859-1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "http default client do")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read responses")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("bad status code: %d, %s", resp.StatusCode, string(b))
	}
	slackresp := &slack.Response{}
	if err := json.Unmarshal(b, slackresp); err != nil {
		return errors.Wrap(err, "unmarshal slack resp")
	}
	if slackresp.Warning != "" {
		logrus.Warning(slackresp.Warning)
	}
	if !slackresp.OK {
		logrus.Error(string(b))
		return errors.Wrap(errors.New(slackresp.Error), "slack response")
	}
	return nil
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

type listUsersResponse struct {
	OK      bool          `json:"ok"`
	Members []*slack.User `json:"members"`
	Error   string        `json:"error"`
}

// ListUsers list users in workspace
func (c *Client) ListUsers() ([]*slack.User, error) {
	data, err := newAPIRequest(slack.UsersList).
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "slack api request")
	}
	listResponse := &listUsersResponse{}
	if err := json.Unmarshal(data, &listResponse); err != nil {
		return nil, errors.Wrap(err, "unmarshal request: %s")
	}
	if !listResponse.OK {
		return nil, errors.New(listResponse.Error)
	}
	return listResponse.Members, nil
}

type listChannelResponse struct {
	OK       bool             `json:"ok"`
	Channels []*slack.Channel `json:"channels"`
	Error    string           `json:"error"`
}

// ListChannels list channels in workspace
func (c *Client) ListChannels() ([]*slack.Channel, error) {
	data, err := newAPIRequest(slack.ChannelsList).
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "slack api request")
	}
	resp := &listChannelResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Channels, nil
}

type listMessagesResponse struct {
	OK       bool                    `json:"ok"`
	Error    string                  `json:"error"`
	Messages []*slack.HistoryMessage `json:"messages"`
	HasMore  bool                    `json:"has_more"`
	PinCount int                     `json:"pin_count"`
}

// ListMessages lists messages in channel
func (c *Client) ListMessages(channelID string) ([]*slack.HistoryMessage, error) {
	data, err := newAPIRequest(slack.ConversationHistory).
		addParam("token", c.BotToken).
		addParam("channel", channelID).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api requests")
	}
	resp := &listMessagesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Messages, nil
}

// ListTodayMessages lists messages from today
func (c *Client) ListTodayMessages(channelID string) ([]*slack.HistoryMessage, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	data, err := newAPIRequest(slack.ConversationHistory).
		addParam("token", c.BotToken).
		addParam("channel", channelID).
		addParam("oldest", fmt.Sprint(today.Unix())).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api requests")
	}
	resp := &listMessagesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Messages, nil
}

type listDirectMessageChannelsReponse struct {
	OK       bool             `json:"ok"`
	Error    string           `json:"error"`
	Channels []*slack.Channel `json:"ims"`
}

// ListDirectMessageChannels lists direct message channels
func (c *Client) ListDirectMessageChannels() ([]*slack.Channel, error) {
	data, err := newAPIRequest(slack.IMList).
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api request")
	}
	resp := &listDirectMessageChannelsReponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Channels, nil
}

type deleteMessageResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// DeleteBotMessage deletes a bot message
func (c *Client) DeleteBotMessage(channelID string, messageTS string) error {
	data, err := newAPIRequest(slack.ChatDelete).
		addParam("channel", channelID).
		addParam("ts", messageTS).
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return errors.Wrap(err, "api request")
	}
	resp := &deleteMessageResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return errors.New(resp.Error)
	}
	return nil
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

type userInfoResponse struct {
	OK    bool        `json:"ok"`
	Error string      `json:"error"`
	User  *slack.User `json:"user"`
}

// GetUserByID searches for a user by iD
func (c *Client) GetUserByID(userID string) (*slack.User, error) {
	data, err := newAPIRequest(slack.UsersInfo).
		addParam("user", userID).
		addParam("include_local", "true").
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api request")
	}
	resp := &userInfoResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.User, nil
}
