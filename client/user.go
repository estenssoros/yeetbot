package client

type User struct {
	Name    string
	SlackID string `yaml:"slackID,omitempty"`
}
