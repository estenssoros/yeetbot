package slack

type Attachment struct {
	Title          string    `json:"title,omitempty"`
	Text           string    `json:"text,omitempty" yaml:"text,omitempty"`
	Fallback       string    `json:"fallback,omitempty" yaml:"fallback,omitempty"`
	CallbackID     string    `json:"callback_id,omitempty" yaml:"callback_id,omitempty"`
	Color          string    `json:"color,omitempty" yaml:"color,omitempty"`
	Actions        []*Action `json:"actions,omitempty" yaml:"actions,omitempty"`
	AttachmentType string    `json:"attachment_type,omitempty" yaml:"attachment_type,omitempty"`
}

func (a *Attachment) AddAction(action *Action) {
	if action.Type == "" {
		action.Type = "button"
	}
	a.Actions = append(a.Actions, action)
}

