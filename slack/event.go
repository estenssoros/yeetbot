package slack

import "encoding/json"

type EventRequest struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		ClientMsgID string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`
		Channel     string `json:"channel"`
		EventTs     int    `json:"event_ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	Type        string   `json:"type"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
	Challenge   string   `json:"challenge"`
}

func (r EventRequest) String() string {
	ju, _ := json.MarshalIndent(r, "", " ")
	return string(ju)
}
