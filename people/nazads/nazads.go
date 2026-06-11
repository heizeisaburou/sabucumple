package nazads

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v5"
)

//go:embed index.gohtml fionne1.png fionne2.png
var fs embed.FS

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "nazads"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("/fionne1.png", serveImage("fionne1.png"))
	g.GET("/fionne2.png", serveImage("fionne2.png"))
}

func home(c *echo.Context) error {
	data, _ := fs.ReadFile("index.gohtml")
	return c.HTML(http.StatusOK, string(data))
}

func serveImage(name string) func(*echo.Context) error {
	return func(c *echo.Context) error {
		data, err := fs.ReadFile(name)
		if err != nil {
			return c.String(http.StatusNotFound, "not found")
		}
		return c.Blob(http.StatusOK, "image/png", data)
	}
}
