// Package juudini es la zona de juudini en el cumple de saburou.
package juudini

import (
	"bytes"
	_ "embed"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

//go:embed index.html
var indexHTML string

//go:embed saburou.mp3
var saburouMP3 []byte

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "juudini"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("", home)
	g.GET("/saburou.mp3", serveAudio)
}

func home(c *echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}

func serveAudio(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "audio/mpeg")
	http.ServeContent(c.Response(), c.Request(), "saburou.mp3", time.Time{}, bytes.NewReader(saburouMP3))
	return nil
}
