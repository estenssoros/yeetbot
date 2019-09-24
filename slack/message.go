package slack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

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
	SubType     string        `json:"subtype,omitempty"`
	Text        string        `json:"text"`
	User        string        `json:"user"`
	Ts          string        `json:"ts"`
	Team        string        `json:"team"`
	BotID       string        `json:"bot_id,omitempty"`
	BotLink     string        `json:"bot_link,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

func (m HistoryMessage) String() string {
	ju, _ := json.MarshalIndent(m, "", " ")
	return string(ju)
}

type listMessagesResponse struct {
	OK       bool              `json:"ok"`
	Error    string            `json:"error"`
	Messages []*HistoryMessage `json:"messages"`
	HasMore  bool              `json:"has_more"`
	PinCount int               `json:"pin_count"`
}

func (a *API) ListMessages(channelID string) ([]*HistoryMessage, error) {
	data, err := a.newRequest(ConversationHistory).
		addParam("token", a.botToken).
		addParam("channel", channelID).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api requests")
	}
	resp := &listMessagesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Messages, nil
}

// ListTodayMessages lists messages from today
func (a *API) ListTodayMessages(channelID string) ([]*HistoryMessage, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	data, err := a.newRequest(ConversationHistory).
		addParam("token", a.botToken).
		addParam("channel", channelID).
		addParam("oldest", fmt.Sprint(today.Unix())).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api requests")
	}
	resp := &listMessagesResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Messages, nil
}

type listDirectMessageChannelsReponse struct {
	OK       bool       `json:"ok"`
	Error    string     `json:"error"`
	Channels []*Channel `json:"ims"`
}

// ListDirectMessageChannels lists direct message channels
func (a *API) ListDirectMessageChannels() ([]*Channel, error) {
	data, err := a.newRequest(IMList).
		addParam("token", a.botToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api request")
	}
	resp := &listDirectMessageChannelsReponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Channels, nil
}

type deleteMessageResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// DeleteBotMessage deletes a bot message
func (a *API) DeleteBotMessage(channelID string, messageTS string) error {
	data, err := a.newRequest(ChatDelete).
		addParam("channel", channelID).
		addParam("ts", messageTS).
		addParam("token", a.botToken).
		Get()
	if err != nil {
		return errors.Wrap(err, "api request")
	}
	resp := &deleteMessageResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return errors.New(resp.Error)
	}
	return nil
}
