package slack

type Action struct {
	Name         string         `json:"name,omitempty" yaml:"name,omitempty"`
	Text         string         `json:"text,omitempty" yaml:"text,omitempty"`
	Type         string         `json:"type,omitempty" yaml:"type,omitempty"`
	Value        string         `json:"value,omitempty" yaml:"value,omitempty"`
	Confirm      *Confirm       `json:"confirm,omitempty" yaml:"confirm,omitempty"`
	Style        string         `json:"style,omitempty" yaml:"style,omitempty"`
	Options      []*Option      `json:"options,omitempty" yaml:"options,omitempty"`
	OptionGroups []*OptionGroup `json:"option_groups,omitempty" yaml:"option_groups,omitempty"`
	DataSource   string         `json:"data_source,omitempty" yaml:"data_source,omitempty"`
}

type ActionURL struct {
	Type            string
	Actions         []*Action
	CallbackID      string
	Team            *Team
	Channel         *Channel
	User            *User
	ActionTS        string
	MessageTS       string
	Token           string
	OriginalMessage *Message
	ResponseURL     string
}
