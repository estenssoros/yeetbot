package client

import (
	"encoding/json"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
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

// FindUserReport selects closest previous report to current time
func (c *Client) FindReportByUser(user *slack.User, userReports map[string][]*Report) (*Report, error) {
	closestTime := struct {
		index int
		time  int64
	}{}
	now := time.Now().Unix()
	if len(userReports[user.RealName]) == 0 {
		return nil, errors.New("No reports found")
	}
	for i, report := range userReports[user.RealName] {
		t, err := time.Parse(time.Kitchen, report.Schedule.Time)
		if err != nil {
			return nil, err
		}
		if t.Unix() < now && t.Unix() > closestTime.time {
			closestTime.index = i
			closestTime.time = t.Unix()
		}
	}
	return userReports[user.RealName][closestTime.index], nil
}

// InitiateReport intiates a new report for a user
func (c *Client) InitiateReport(user *slack.User) error {
	// TODO this
	return nil
}

// IsReportComplete checks to see if the user has completed all questions
func (c *Client) IsReportComplete(user *slack.User) bool {
	// TODO this
	return true
}

// CompleteReport sends the users report to slack
func (c *Client) CompleteReport(user *slack.User) error {
	// TODO this
	return nil
}
