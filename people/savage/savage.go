package savage

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
  templates: template.Must(template.ParseGlob("people/savage/templates/*.html")),
}

const Name = "savage"

func (Module) Endpoint() string {
  return Name
}

func (Module) Register(g *echo.Group) {
  g.GET("/", home)
  g.GET("", home)
  g.Static("/assets", "people/savage/assets")
}

func home(c *echo.Context) error {
  data := map[string]any{"Name": Name}
  return t.Render(c, c.Response(), "main.html", data)
}
