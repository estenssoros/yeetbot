package client

import (
	"encoding/json"
)

type Config struct {
	UserToken string    `json:"user_token"`
	BotToken  string    `json:"bot_token"`
	Reports   []*Report `json:"reports"`
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

func (c *Config) NewClient(report *Report) *Client {
	return &Client{
		UserToken: c.UserToken,
		BotToken:  c.BotToken,
		Report:    report,
	}
}
