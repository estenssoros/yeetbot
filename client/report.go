package client

import "encoding/json"

type Report struct {
	Name         string      `json:"name"`
	Channel      string      `json:"channel"`
	Users        []*User     `json:"users"`
	Schedule     *Schedule   `json:"schedule"`
	IntroMessage string      `json:"intro_message"`
	Questions    []*Question `json:"questions"`
}

func (r Report) String() string {
	ju, _ := json.MarshalIndent(r, "", " ")
	return string(ju)
}
