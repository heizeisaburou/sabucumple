package ejemplo

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "ejemplo"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("/otro", otro)
}

func home(c *echo.Context) error {
	return c.HTML(http.StatusOK, `
		<h1>Zona de ejemplo</h1>
		<p>Aquí ejemplo puede hacer lo que quiera.</p>
		<a href="/ejemplo/otro">Ver ejemplo</a>
	`)
}

func otro(c *echo.Context) error {
	return c.HTML(http.StatusOK, `
		<h1>SubEndpoint de ejemplo</h1>
		<p>blablabla.</p>
		<a href="/ejemplo/">Volver</a>
	`)
}
