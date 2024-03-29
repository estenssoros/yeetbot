package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/estenssoros/yeetbot/views"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

func logError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
}

func main() {
	engine := echo.New()
	engine.Use(middleware.Recover())
	engine.Use(client.Middleware)
	engine.Use(logError)
	engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] - ${time_rfc3339} |${status}| ${latency_human} | ${host} | ${method} ${uri}\n",
	}))
	engine.POST("/event", views.EventHandler)
	engine.POST("/interact", views.InteractHandler)
	for _, r := range engine.Routes() {
		fmt.Println(r.Method, r.Path, r.Name)
	}
	engine.Start(":3000")
}
