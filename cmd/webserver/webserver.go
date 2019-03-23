package webserver

import (
	"github.com/joostvdg/cmg/pkg/webserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

func StartWebserver() {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}


	// Echo instance
	e := echo.New()
	e.Logger.Printf("Starting server on port %s\n", port)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/api/map", webserver.GetMap)
	e.GET("/api/legend", webserver.GetMapLegend)

	// Start server
	e.Logger.Fatal(e.Start(":"+port))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
