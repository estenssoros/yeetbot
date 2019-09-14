package client

// Question a question to ask a user
// Color is the attachment color
// If options are given, the question will have a drop down
type Question struct {
	Text    string   `yaml:"text"`
	Color   string   `yaml:"color"`
	Options []string `yaml:"options,omitempty"`
}

// PostFirstQuestion sends the first question to the user
func (c *Client) PostFirstQuestion() error {
	return nil
}

// PostNextQuestion sends the next question to the user
func (c *Client) PostNextQuestion() error {
	// TODO check to see if we are on first question and initiate report
	return nil
}

// GetLastQuestion gets the last question asked by yeetbot
func (c *Client) GetRecentQuestion() (*Question, error) {
	// TODO this
	return nil, nil
}
