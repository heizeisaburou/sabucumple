// Package midos es la zona de MidOS en el cumple de saburou.
//
// Más que un saludo, es un pedacito del organismo: un módulo autocontenido
// que no toca el estado global de la app (su propio renderer, sus propios
// contadores) y que expresa cómo trabaja MidOS — rieles que guían, nodos que
// deben estar todos en coherencia, certeza que se calcula, no se adivina.
//
// Patrones aplicados (los mismos que MidOS usa en producción):
//   - html/template + go:embed: el HTML vive en assets/, no concatenado a mano.
//   - Render a buffer local: cero mutación del echo.Renderer global → no pisa
//     la zona de nadie más (autocontención = coherencia).
//   - atomic counters: métricas lock-free en el hot path.
//   - fail-closed: si un template falla, se responde el error, nunca a medias.
package midos

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v5"
)

// assets embebe los templates en el binario. Cero archivos sueltos en deploy:
// el `go build` produce un único ejecutable con todo adentro.
//
//go:embed assets/*.gohtml
var assets embed.FS

// tmpl se parsea una sola vez al cargar el paquete. template.Must hace
// fail-fast: si un template está roto, el programa no arranca (mejor que
// descubrirlo en el primer request).
var tmpl = template.Must(template.ParseFS(assets, "assets/*.gohtml"))

// startedAt marca el arranque del paquete — base del uptime en /salud.
var startedAt = time.Now()

// requests cuenta las visitas a la zona sin locks (atomic en hot path).
var requests atomic.Int64

// Module es la zona de MidOS. Value receiver a propósito: la interfaz
// module.Module no necesita punteros (un Module no tiene estado mutable propio;
// el estado vive en los counters a nivel paquete).
type Module struct{}

// New construye la zona. Firma exacta que main.go espera: midos.New().
func New() Module {
	return Module{}
}

// Endpoint es la ruta raíz de la zona: /midos.
func (Module) Endpoint() string {
	return "midos"
}

// Register cablea las rutas de la zona dentro de su grupo.
//
//	GET /midos            → el regalo: qué es MidOS, en una página
//	GET /midos/coherencia → demo viva del AND de nodos (un gate en miniatura)
//	GET /midos/salud      → health-check estilo MidOS (uptime + contadores)
func (Module) Register(g *echo.Group) {
	g.GET("/", home)
	g.GET("/coherencia", coherencia)
	g.GET("/salud", salud)
}

// render ejecuta un template contra un buffer local y recién entonces escribe
// la respuesta. Si el template falla, no se emitió ni un byte: el handler puede
// devolver 500 limpio en lugar de un HTML cortado a la mitad (fail-closed).
func render(c *echo.Context, name string, data any) error {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

// home sirve el regalo: una página que explica qué es MidOS sin marketing —
// la metáfora tranvía/rieles/cable y los nodos de coherencia.
func home(c *echo.Context) error {
	requests.Add(1)
	return render(c, "home.gohtml", nil)
}

// nodo es un eslabón de coherencia. En MidOS la coherencia no es un score
// difuso: es el AND de nodos (doctrina · doc · diagrama · test · código …).
// Si uno solo está en false, el todo está incoherente.
type nodo struct {
	Nombre string
	OK     bool
}

// coherenciaData alimenta el template de la demo.
type coherenciaData struct {
	Nodos     []nodo
	Veredicto string // "COHERENTE" si todos OK, si no "INCOHERENTE"
	Roto      string // nombre del primer nodo en false (vacío si todos OK)
}

// evaluar implementa el corazón de un gate MidOS: AND de nodos, corto-circuito
// al primer false. Los gates reales del organismo tienen exactamente esta
// forma — un gate es, en el fondo, esto.
func evaluar(nodos []nodo) (string, string) {
	for _, n := range nodos {
		if !n.OK {
			return "INCOHERENTE", n.Nombre
		}
	}
	return "COHERENTE", ""
}

// coherencia es un gate en miniatura, ejecutable desde el navegador.
// Cada nodo se activa con ?nodo=true (default true). Apagá uno con
// ?doctrina=false y mirá cómo el veredicto colapsa al primer false: así razona
// un riel de MidOS.
func coherencia(c *echo.Context) error {
	requests.Add(1)

	// Default: todos los nodos en coherencia. El visitante los apaga por query.
	flag := func(name string) bool {
		return c.QueryParam(name) != "false"
	}
	nodos := []nodo{
		{"doctrina", flag("doctrina")},
		{"doc", flag("doc")},
		{"diagrama", flag("diagrama")},
		{"test", flag("test")},
		{"código", flag("codigo")},
		{"rieles", flag("rieles")},
	}
	veredicto, roto := evaluar(nodos)

	return render(c, "coherencia.gohtml", coherenciaData{
		Nodos:     nodos,
		Veredicto: veredicto,
		Roto:      roto,
	})
}

// salud responde un health-check estilo MidOS gateway: JSON con uptime y
// contadores atómicos. Verdict ALLOW fijo — la zona no tiene backend que pueda
// caerse; si respondés esto, estás viva.
func salud(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"verdict":  "ALLOW",
		"zona":     "midos",
		"uptime_s": int(time.Since(startedAt).Seconds()),
		"requests": requests.Load(),
		"para":     "三郎",
	})
}
