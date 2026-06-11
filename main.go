package main

import (
  "fmt"
  "net/http"
  "strings"

  "github.com/labstack/echo/v5"

  "github.com/heizeisaburou/sabucumple/module"
  "github.com/heizeisaburou/sabucumple/people/ejemplo"
  "github.com/heizeisaburou/sabucumple/people/midos"
  "github.com/heizeisaburou/sabucumple/people/savage"
)

func main() {
  e := echo.New()

  modules := []module.Module{
    ejemplo.New(),
    midos.New(),
    savage.New(),
  }


  e.GET("/", func(c *echo.Context) error {
    var html strings.Builder

    html.WriteString(`
      <h1>Cumple 🎂</h1>
      <ul>
    `)

    for _, m := range modules {
      html.WriteString(`<li><a href="/` + m.Endpoint() + `/">` + m.Endpoint() + `</a></li>`)
    }

    html.WriteString(`</ul>`)

    return c.HTML(http.StatusOK, html.String())
  })

  for _, m := range modules {
    g := e.Group("/" + m.Endpoint())

    // Rutas propias de cada persona.
    m.Register(g)
  }

  e.Static("/static", "static")

  fmt.Println("Servidor escuchando en http://localhost:8080")
  e.Start(":8080")
}
