package client

import "github.com/estenssoros/yeetbot/slack"

// User naive user data structure
type User struct {
	Name    string
	SlackID string `yaml:"slackID,omitempty"`
}

func (c *Client) HasUser(user *slack.User) bool {
	for _, u := range c.Users {
		if u.Name == user.Name {
			return true
		}
	}
	return false
}
