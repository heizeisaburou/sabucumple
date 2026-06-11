package chavsi

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed index.html
var indexHTML string

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "chavsi"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", chavsi)
}

func chavsi(c *echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}