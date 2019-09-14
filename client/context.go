package client

import (
	"fmt"
	"log"

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

func (c *Context) populateUserReports() {
	for _, report := range c.Config.Reports {
		for _, user := range report.Users {
			c.UserReports[user.ID] = append(c.UserReports[user.Name], report)
		}
	}
}

func (c *Context) NewClientFromUser(userID string) (*Client, error) {
	reports, ok := c.UserReports[userID]
	if !ok {
		fmt.Println(c.Config)
		return nil, errors.Errorf("could not locate user: %s", userID)
	}
	// TODO how do we find the right report?
	return c.Config.NewClient(reports[0]), nil
}
