package client

import (
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
)

var (
	userTimeZone = "user-local"
	timePattern  = "15:04"
)

// Schedule records period, dow, timezone, and excluded holidays
type Schedule struct {
	Period          string `yaml:"period"`
	Mon             bool   `yaml:"mon"`
	Tues            bool   `yaml:"tues"`
	Wed             bool   `yaml:"wed"`
	Thurs           bool   `yaml:"thurs"`
	Fri             bool   `yaml:"fri"`
	Sat             bool   `yaml:"sat"`
	Sun             bool   `yaml:"sun"`
	Time            string `yaml:"time"`
	TimeZone        string `yaml:"timeZone"`
	ExcludeHolidays string `yaml:"excludeHolidays,omitempty"`
}

func (s *Schedule) parseTime() (time.Time, error) {
	return time.Parse(timePattern, s.Time)
}

// TodayTime converts schedule time to today's time in UTC
func (s *Schedule) TodayTime() (time.Time, error) {
	t, err := s.parseTime()
	if err != nil {
		return t, errors.Wrap(err, "schedule parse time")
	}
	today := time.Now()
	return time.Date(today.Year(), today.Month(), today.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC), nil
}

// UserTimeZone gets today time and add a users slacks tz offset
func (s *Schedule) UserTimeZone(user *slack.User) (time.Time, error) {
	t, err := s.TodayTime()
	if err != nil {
		return t, errors.Wrap(err, "schedule parse time")
	}
	return t.Add(time.Second * time.Duration(user.TZOffset)), nil
}
