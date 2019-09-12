package client

import (
	"time"

	"github.com/pkg/errors"
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

// TodayTime converts schedule time to today's time
// TODO if report says user timezone, how do we handle
func (s *Schedule) TodayTime() (time.Time, error) {
	t, err := time.Parse("15:04", s.Time)
	if err != nil {
		return t, err
	}
	loc, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return t, errors.Wrapf(err, "time load location %s", s.TimeZone)
	}
	today := time.Now().UTC()
	return time.Date(today.Year(), today.Month(), today.Day(), t.Hour(), t.Minute(), 0, 0, loc), nil
}
