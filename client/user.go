package client

// User naive user data structure
type User struct {
	Name    string
	SlackID string `yaml:"slackID,omitempty"`
}
