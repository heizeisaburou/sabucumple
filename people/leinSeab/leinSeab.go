package leinSeab

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed assets/*
var assets embed.FS

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "leinSeab"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", index)
	g.GET("", index)
	g.GET("/style.css", serveCSS)
	g.GET("/main.js", serveJS)
}

func index(c *echo.Context) error {
	htmlBytes, err := assets.ReadFile("assets/index.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, string(htmlBytes))
}

func serveCSS(c *echo.Context) error {
	cssBytes, err := assets.ReadFile("assets/style.css")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Blob(http.StatusOK, "text/css; charset=utf-8", cssBytes)
}

func serveJS(c *echo.Context) error {
	jsBytes, err := assets.ReadFile("assets/main.js")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Blob(http.StatusOK, "application/javascript; charset=utf-8", jsBytes)
}
