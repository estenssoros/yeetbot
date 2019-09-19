package client

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
	"gopkg.in/yaml.v2"
)

var (
	yeetENV             = "YEETBOT_CONFIG"
	defaultMeetingIndex = "yeetmeet"
	defaultReportIndex  = "yeetreport"
)

// Config all info for a yeetbot config
type Config struct {
	Team      string            `json:"team" yaml:"team"`
	UserToken string            `json:"userToken" yaml:"userToken"`
	BotToken  string            `json:"botToken" yaml:"botToken"`
	YeetUser  string            `json:"yeetUserID" yaml:"yeetUser"`
	Debug     bool              `json:"debug" yaml:"debug"`
	Meetings  []*models.Meeting `json:"meetings" yaml:"meetings"`
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// NewClient creates a report client from a config
func (c *Config) NewClient() *Client {
	s := slack.New(c.BotToken)
	s.SetVerbose(c.Debug)
	client := &Client{
		UserToken:    c.UserToken,
		YeetUser:     c.YeetUser,
		BotToken:     c.BotToken,
		Config:       c,
		reportIndex:  defaultReportIndex,
		meetingIndex: defaultMeetingIndex,
		slack:        s,
		elastic:      elasticsvc.New(context.Background()),
	}
	return client
}

// ConfigFromReader creates a config from a reader
func ConfigFromReader(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read reader")
	}
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal")
	}
	return config, nil
}

// ConfigFromEnv loads config from an environment variable
func ConfigFromEnv() (*Config, error) {
	path := os.Getenv(yeetENV)
	if path == "" {
		return nil, errors.New("missing " + yeetENV)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer f.Close()
	return ConfigFromReader(f)
}
