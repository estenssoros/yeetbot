package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	yeetENV      = "YEETBOT_CONFIG"
	elasticIndex = "yeetbot"
)

// Config all info for a yeetbot config
type Config struct {
	UserToken string    `json:"user_token"`
	BotToken  string    `json:"bot_token"`
	YeetUser  string    `json::""yeet_userid`
	Debug     bool      `json:"debug"`
	Reports   []*Report `json:"reports"`

	ElasticURL string `json:"elastic_url"`
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// NewClient creates a report client from a config
func (c *Config) NewClient(report *Report) *Client {
	client := &Client{
		YeetUser:     c.YeetUser,
		ElasticIndex: elasticIndex,
		UserReports:  map[string][]*Report{},
		UserMap:      map[string]*slack.User{},
		Config:       c,
		Report:       report,
	}
	client.PopulateUserReports()
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
	logrus.Infof("using config stored at: %s", path)
	return ConfigFromReader(f)
}
