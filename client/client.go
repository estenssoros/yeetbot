package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	slackurl        = "https://slack.com/api"
	defaultUsername = "yeetbot"
	defaultIcon     = ":ghost:"
	defaultChannel  = "#general"
	yeetbotBucket   = "yeetbot"
)

type Config struct {
	URL       string  `yaml:"url"`
	UserToken string  `yaml:"user_token"`
	BotToken  string  `yaml:"bot_token"`
	Username  string  `yaml:"username"`
	Icon      string  `yaml:"icon"`
	Channel   string  `yaml:"channel"`
	Debug     bool    `yaml:"debug"`
	Greeting  string  `yaml:"greeting"`
	Steps     []*Step `yaml:"steps"`
}

func ReadConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "readfile")
	}
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal")
	}
	return config, nil
}

type Client struct {
	URL       string       `json:"url"`
	UserToken string       `json:"user_token"`
	BotToken  string       `json:"bot_token"`
	Username  string       `json:"username"`
	Icon      string       `json:"icon"`
	Channel   string       `json:"channel"`
	Debug     bool         `json:"debug"`
	Greeting  string       `json:"greeting"`
	Steps     []*Step      `json:"steps"`
	Team      *models.Team `json:"team"`
}

func (c Client) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

func ClientFromConfig(fileName string) (*Client, error) {
	config, err := ReadConfig(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}
	return New(config), nil
}

func (c *Client) AddTeamFromFile(fileName string) error {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrap(err, "readfile")
	}
	team := &models.Team{}
	if err := yaml.Unmarshal(d, team); err != nil {
		return errors.Wrap(err, "yaml unmarshal")
	}
	c.Team = team
	return nil
}

func New(config *Config) *Client {
	return &Client{
		URL:       config.URL,
		UserToken: config.UserToken,
		BotToken:  config.BotToken,
		Username:  config.Username,
		Icon:      config.Icon,
		Channel:   config.Channel,
		Debug:     config.Debug,
		Greeting:  config.Greeting,
		Steps:     config.Steps,
	}
}

func NewAWS() (*Client, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))
	if err != nil {
		return nil, errors.Wrap(err, "session new")
	}
	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}
	_, err = downloader.Download(buff, &s3.GetObjectInput{
		Bucket: &yeetbotBucket,
		Key:    aws.String("config/config.yaml"),
	})
	if err != nil {
		return nil, errors.Wrap(err, "download")
	}
	config := &Config{}
	if err := yaml.Unmarshal(buff.Bytes(), config); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal")
	}
	return New(config), nil
}

func (c *Client) applyMessageDefaults(msg *slack.Message) {
	if msg.Username == "" {
		msg.Username = defaultUsername
	}
	if msg.Icon == "" {
		msg.Icon = defaultIcon
	}
	if msg.Channel == "" {
		msg.Channel = defaultChannel
	}
}

func (c *Client) SendMessage(msg *slack.Message) error {
	c.applyMessageDefaults(msg)
	u, err := url.Parse(slackurl + "/" + slack.ChatPostMessage)
	if err != nil {
		return errors.Wrap(err, "url parse")
	}
	if err := c.postRequest(u.String(), msg); err != nil {
		return errors.Wrap(err, "client post request")
	}
	return nil
}

func (c *Client) postRequest(url string, v interface{}) error {
	if c.BotToken == "" {
		return errors.New("missing bot token")
	}
	if c.UserToken == "" {
		return errors.New("missing user token")
	}
	if c.Debug {
		logrus.Infof("%s %s", http.MethodPost, url)
	}
	ju, _ := json.Marshal(v)
	if c.Debug {
		logrus.Info(string(ju))
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(ju))
	if err != nil {
		return errors.Wrap(err, "http new request")
	}
	{
		req.Header.Add("Authorization", "Bearer "+c.BotToken)
		req.Header.Add("Content-Type", "application/json;charset=iso-8859-1")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "http default client do")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read responses")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("bad status code: %d, %s", resp.StatusCode, string(data))
	}
	slackresp := &slack.Response{}
	if err := json.Unmarshal(data, slackresp); err != nil {
		return errors.Wrap(err, "unmarshal slack resp")
	}
	if slackresp.Warning != "" {
		logrus.Warning(slackresp.Warning)
	}
	if !slackresp.OK {
		logrus.Error(slackresp)
		return errors.Wrap(errors.New(slackresp.Error), "slack response")
	}
	return nil
}

func (c *Client) GenericMessage(u *models.User, text string) error {
	msg := &slack.Message{
		Text:    text,
		Channel: "@" + u.Name,
		AsUser:  true,
	}
	return errors.Wrap(c.SendMessage(msg), "client send message")
}

func (c *Client) SendGreeting(user *models.User) error {
	if c.BotToken == "" {
		return errors.New("missing bot token")
	}
	if c.UserToken == "" {
		return errors.New("missing user token")
	}
	text, err := user.Template(c.Greeting)
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

func (c *Client) Run(u *models.User) error {
	if c.Team == nil {
		return errors.New("no team configured")
	}
	for i, s := range c.Steps {
		if err := c.GenericMessage(u, s.Text); err != nil {
			return errors.Wrapf(err, "generic message step: %d", i)
		}
	}
	return nil
}

type listUsersResponse struct {
	OK      bool           `json:"ok"`
	Members []*models.User `json:"members"`
}

func (c *Client) ListUsers() ([]*models.User, error) {
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
	return listResponse.Members, nil
}

type listChannelResponse struct {
	OK       bool             `json:"ok"`
	Channels []*slack.Channel `json:"channels"`
}

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
	return resp.Channels, nil
}

func (c *Client) GetUserFromRequest(req *models.EventRequest) (*models.User, error) {
	users, err := c.ListUsers()
	if err != nil {
		return nil, errors.Wrap(err, "client list users")
	}
	for _, u := range users {
		if u.ID == req.Event.User {
			return u, nil
		}
	}
	return nil, errors.Errorf("could not locate user %s", req.Event.User)
}

func (c *Client) ListConversations() (interface{}, error) {
	data, err := newAPIRequest(slack.ConversationsList).
		addParam("token", c.BotToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "slack api request")
	}
	fmt.Println(string(data))
	return nil, errors.New("not implemented")
}

func (c *Client) GetLastMessageFromUser(user *models.User) (*slack.Message, error) {
	conversations, err := c.ListConversations()
	if err != nil {
		return nil, errors.Wrap(err, "client list conversations")
	}
	fmt.Println(conversations)
	return nil, nil
}

type listMessagesResponse struct {
	OK       bool                    `json:"ok"`
	Error    string                  `json:"error"`
	Messages []*slack.HistoryMessage `json:"messages"`
	HasMore  bool                    `json:"has_more"`
	PinCount int                     `json:"pin_count"`
}

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

type listDirectMessageChannelsReponse struct {
	OK       bool             `json:"ok"`
	Error    string           `json:"error"`
	Channels []*slack.Channel `json:"ims"`
}

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

func (c *Client) DeleteUserMessage(channelID string, messageTS string) error {
	data, err := newAPIRequest(slack.ChatDelete).
		addParam("channel", channelID).
		addParam("ts", messageTS).
		addParam("token", c.UserToken).
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
