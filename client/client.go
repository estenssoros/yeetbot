package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

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
)

type Config struct {
	URL      string
	Token    string
	Username string
	Icon     string
	Channel  string
	Debug    bool
	Greeting string
	Steps    []*Step
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
	URL      string       `json:"url"`
	Token    string       `json:"token"`
	Username string       `json:"username"`
	Icon     string       `json:"icon"`
	Channel  string       `json:"channel"`
	Debug    bool         `json:"debug"`
	Greeting string       `json:"greeting"`
	Steps    []*Step      `json:"steps"`
	Team     *models.Team `json:"team"`
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
		URL:      config.URL,
		Token:    config.Token,
		Username: config.Username,
		Icon:     config.Icon,
		Channel:  config.Channel,
		Debug:    config.Debug,
		Greeting: config.Greeting,
		Steps:    config.Steps,
	}
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
	if c.Debug {
		logrus.Info(http.MethodPost, url)
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
		req.Header.Add("Authorization", "Bearer "+c.Token)
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
	if c.Token == "" {
		return errors.New("missing token")
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
				Text:  "When you are ready, please answer the following question:",
				Color: "#3AA3E3",
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
								Text:  "I'm on vacation 🏝–",
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
