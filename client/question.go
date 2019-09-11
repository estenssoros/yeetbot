package client

type Question struct {
	Text    string   `yaml:"text"`
	Color   string   `yaml:"color"`
	Options []string `yaml:"options,omitempty"`
}
