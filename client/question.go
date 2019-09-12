package client

// Question a question to ask a user
// Color is the attachment color
// If options are given, the question will have a drop down
// TODO: this ^
type Question struct {
	Text    string   `yaml:"text"`
	Color   string   `yaml:"color"`
	Options []string `yaml:"options,omitempty"`
}
