package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	yeetENV = "YEETBOT_CONFIG"
)

// Config all info for a yeetbot config
type Config struct {
	UserToken  string    `json:"user_token"`
	BotToken   string    `json:"bot_token"`
	ElasticURL string    `json:"elastic_url"`
	Debug      bool      `json:"debug"`
	Reports    []*Report `json:"reports"`
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// NewClient creates a report client from a config
func (c *Config) NewClient(report *Report) *Client {
	return &Client{
		UserToken: c.UserToken,
		BotToken:  c.BotToken,
		Debug:     c.Debug,
		Report:    report,
	}
}

func (c *Config) NewClientFromChannel(channel string) (*Client, error) {
	for _, report := range c.Reports {
		if report.Channel == channel {
			return c.NewClient(report), nil
		}
	}
	return nil, errors.New("unable to find report")
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
