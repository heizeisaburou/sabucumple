package saburou

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed index.gohtml
var indexHTML string

//go:embed assets/styles.css
var stylesCSS []byte

//go:embed assets/flonne_blue_curr.png assets/flonne_red_curr.png
var assets embed.FS

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "saburou"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("", home)
	g.GET("/styles.css", serveCSS)
	g.GET("/assets/flonne_blue_curr.png", serveImage("assets/flonne_blue_curr.png"))
	g.GET("/assets/flonne_red_curr.png", serveImage("assets/flonne_red_curr.png"))
}

func home(c *echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}

func serveCSS(c *echo.Context) error {
	return c.Blob(http.StatusOK, "text/css; charset=utf-8", stylesCSS)
}

func serveImage(path string) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		data, err := assets.ReadFile(path)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.Blob(http.StatusOK, "image/png", data)
	}
}
