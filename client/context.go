package client

import (
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Context wrapper around echo's context
type Context struct {
	echo.Context
	Config      *Config
	UserReports map[string][]*Report
}

// Middleware to wrap echo's context with Context
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	config, err := ConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}
	return func(c echo.Context) error {
		cc := &Context{Context: c, Config: config, UserReports: map[string][]*Report{}}
		cc.populateUserReports()
		return next(cc)
	}
}

// FindUserReport selects closest previous report to current time
func (c *Context) FindUserReport(username string) (*Report, error) {
	closestTime := struct {
		index int
		time  int64
	}{}
	now := time.Now().Unix()
	if len(c.UserReports[username]) == 0 {
		return nil, errors.New("No reports found")
	}
	for i, report := range c.UserReports[username] {
		t, err := time.Parse(time.Kitchen, report.Schedule.Time)
		if err != nil {
			return nil, err
		}
		if t.Unix() < now && t.Unix() > closestTime.time {
			closestTime.index = i
			closestTime.time = t.Unix()
		}
	}
	return c.UserReports[username][closestTime.index], nil
}

func (c *Context) populateUserReports() {
	for _, report := range c.Config.Reports {
		for _, user := range report.Users {
			c.UserReports[user.Name] = append(c.UserReports[user.Name], report)
		}
	}
}
