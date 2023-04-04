package main

import (
	"github.com/a-paudel/fileshare/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func startServer() {
	app := echo.New()
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${status}\t${method}\t${uri}\t${latency_human}\n",
	}))
	app.Use(middleware.Recover())

	RegisterTemplates(app)
	routes.RegisterFileRoutes(app)

	app.Start(":8000")
}
