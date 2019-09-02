package slack

type Option struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type OptionGroup struct {
	Text    string    `json:"text"`
	Options []*Option `json:"options"`
}

type OptionsLoadURL struct {
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	CallBackID   string   `json:"call_back_id"`
	Type         string   `json:"type"`
	Team         *Team    `json:"team"`
	Channel      *Channel `json:"channel"`
	User         *User    `json:"user"`
	ActionTS     string   `json:"action_ts"`
	AttachmentID string   `json:"attachment_id"`
	Token        string   `json:"token"`
}
