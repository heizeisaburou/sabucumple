package kagliostro

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

type Module struct{}

func New() Module {
	return Module{}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var t = Template{
	templates: template.Must(template.ParseGlob("people/kagliostro/assets/*.gohtml")),
}

func (Module) Endpoint() string {
	return "kagliostro"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("", home)
	g.Static("/assets", "people/kagliostro/assets")
}

func home(c *echo.Context) error {
	return t.Render(c, c.Response(), "home.gohtml", nil)
}
