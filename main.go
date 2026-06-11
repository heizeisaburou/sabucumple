package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v5"

	"github.com/heizeisaburou/sabucumple/module"
	"github.com/heizeisaburou/sabucumple/people/chavsi"
	"github.com/heizeisaburou/sabucumple/people/ejemplo"
	"github.com/heizeisaburou/sabucumple/people/leinSeab"
	"github.com/heizeisaburou/sabucumple/people/midos"
	"github.com/heizeisaburou/sabucumple/people/nazads"
	"github.com/heizeisaburou/sabucumple/people/saburou"
	"github.com/heizeisaburou/sabucumple/people/savage"
	"github.com/heizeisaburou/sabucumple/people/ladythekilla"
  "github.com/heizeisaburou/sabucumple/people/kagliostro"
)

func main() {
	e := echo.New()

	modules := []module.Module{
		ejemplo.New(),
		midos.New(),
		nazads.New(),
		savage.New(),
		chavsi.New(),
		saburou.New(),
		leinSeab.New(),
		ladythekilla.New(),
    kagliostro.New(),
	}

	e.GET("/", saburou.Home)

	for _, m := range modules {
		g := e.Group("/" + m.Endpoint())

		// Rutas propias de cada persona.
		m.Register(g)
	}

	e.Static("/static", "static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Servidor escuchando en :" + port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
