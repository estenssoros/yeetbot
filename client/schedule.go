package client

import (
	"time"

	"github.com/gorhill/cronexpr"
)

var (
	userTimeZone = "user-local"
	timePattern  = "15:04"
)

// Schedule records period, dow, timezone, and excluded holidays
type Schedule struct {
	Cron            string `json:"cron"`
	TimeZone        string `yaml:"timeZone"`
	ExcludeHolidays string `yaml:"excludeHolidays,omitempty"`
}

// NextReportDate figures out next report date
func (s *Schedule) NextReportDate() time.Time {
	return cronexpr.MustParse(s.Cron).Next(time.Now())
}
