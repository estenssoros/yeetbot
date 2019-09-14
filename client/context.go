package client

import (
	"log"

	"github.com/labstack/echo"
)

// Context wrapper around echo's context
type Context struct {
	echo.Context
	Config *Config
}

// Middleware to wrap echo's context with Context
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	config, err := ConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}
	return func(c echo.Context) error {
		cc := &Context{Context: c, Config: config}
		return next(cc)
	}
}
