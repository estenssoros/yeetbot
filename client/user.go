package client

import "github.com/estenssoros/yeetbot/slack"

// User naive user data structure
type User struct {
	Name string
	ID   string `yaml:"id,omitempty"`
}

func (c *Client) HasUser(user *slack.User) bool {
	for _, u := range c.Users {
		if u.Name == user.Name {
			return true
		}
	}
	return false
}

// HasUserStartedReport checks to see if a report has already been started today
func (c *Client) HasUserStartedReport(user *slack.User) bool {
	// TODO this
	return false
}
