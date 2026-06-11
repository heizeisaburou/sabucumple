package saburou

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed index.gohtml
var indexHTML string

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "saburou"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
}

func home(c *echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}
