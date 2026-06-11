package ladythekilla

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed index.html
var indexHTML string

//go:embed assets
var assetsFS embed.FS

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "ladythekilla"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	sub, _ := fs.Sub(assetsFS, "assets")
	g.StaticFS("/assets", sub)
}

func home(c *echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}
