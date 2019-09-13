package client

import (
	"encoding/json"
	"time"
)

// Report a yeetbot report
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

// TodayTime return the schedule report time for today
func (r *Report) TodayTime() (time.Time, error) {
	return r.Schedule.TodayTime()
}
