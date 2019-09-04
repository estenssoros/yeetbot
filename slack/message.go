package slack

import "encoding/json"

type Message struct {
	Text            string        `json:"text,omitempty" yaml:"text,omitempty"`
	AsUser          bool          `json:"as_user,omitempty" yaml:"as_user,omitempty"`
	Username        string        `json:"username,omitempty" yaml:"username,omitempty"`
	Icon            string        `json:"icon_emoji,omitempty" yaml:"icon_emoji,omitempty"`
	Channel         string        `json:"channel,omitempty" yaml:"channel,omitempty"`
	Attachments     []*Attachment `json:"attachments,omitempty" yaml:"attachments,omitempty"`
	ThreadTs        *string       `json:"thread_ts,omitempty" yaml:"thread_ts,omitempty"`
	ResponseType    *string       `json:"response_type,omitempty" yaml:"response_type,omitempty"`
	ReplaceOriginal *bool         `json:"replace_original,omitempty" yaml:"replace_original,omitempty"`
	DeleteOriginal  *bool         `json:"delete_original,omitempty" yaml:"delete_original,omitempty"`
}

func (m *Message) AddAttachment(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

type HistoryMessage struct {
	ClientMsgID string        `json:"client_msg_id,omitempty"`
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	User        string        `json:"user"`
	Ts          string        `json:"ts"`
	Team        string        `json:"team"`
	BotID       string        `json:"bot_id,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

func (m HistoryMessage) String() string {
	ju, _ := json.MarshalIndent(m, "", " ")
	return string(ju)
}
