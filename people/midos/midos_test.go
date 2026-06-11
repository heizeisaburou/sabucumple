package midos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

// newCtx arma un echo.Context de prueba apuntando a `target` (ej "/midos?x=y").
// Helper baseline al estilo MidOS: cada test parte de acá y sólo cambia lo suyo.
func newCtx(target string) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// TestEvaluar cubre el corazón del gate: AND de nodos con corto-circuito.
// Table-driven — un caso por fila, el estado mutado es evidente.
func TestEvaluar(t *testing.T) {
	cases := []struct {
		name      string
		nodos     []nodo
		veredicto string
		roto      string
	}{
		{
			name:      "todos OK → coherente",
			nodos:     []nodo{{"a", true}, {"b", true}, {"c", true}},
			veredicto: "COHERENTE",
			roto:      "",
		},
		{
			name:      "uno roto → incoherente, nombra el nodo",
			nodos:     []nodo{{"a", true}, {"b", false}, {"c", true}},
			veredicto: "INCOHERENTE",
			roto:      "b",
		},
		{
			name:      "corto-circuito al PRIMER false",
			nodos:     []nodo{{"a", false}, {"b", false}},
			veredicto: "INCOHERENTE",
			roto:      "a", // no "b": corta en el primero
		},
		{
			name:      "lista vacía → coherente trivial (AND de nada es true)",
			nodos:     nil,
			veredicto: "COHERENTE",
			roto:      "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ver, roto := evaluar(tc.nodos)
			if ver != tc.veredicto {
				t.Errorf("veredicto: got %q, want %q", ver, tc.veredicto)
			}
			if roto != tc.roto {
				t.Errorf("roto: got %q, want %q", roto, tc.roto)
			}
		})
	}
}

// TestHome verifica que el regalo se renderiza y contiene lo esencial.
func TestHome(t *testing.T) {
	c, rec := newCtx("/midos")
	if err := home(c); err != nil {
		t.Fatalf("home() error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"三郎", "MidOS", "Tranvía", "Rieles", "Cable"} {
		if !strings.Contains(body, want) {
			t.Errorf("home no contiene %q", want)
		}
	}
}

// TestCoherencia_Handler comprueba el render del gate según query params.
func TestCoherencia_Handler(t *testing.T) {
	cases := []struct {
		name       string
		target     string
		wantVer    string
		wantInBody string
	}{
		{"default todos verdes", "/midos/coherencia", "COHERENTE", "COHERENTE"},
		{"doctrina apagada", "/midos/coherencia?doctrina=false", "INCOHERENTE", "doctrina"},
		{"test apagado", "/midos/coherencia?test=false", "INCOHERENTE", "test"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newCtx(tc.target)
			if err := coherencia(c); err != nil {
				t.Fatalf("coherencia() error: %v", err)
			}
			if rec.Code != http.StatusOK {
				t.Fatalf("status: got %d, want 200", rec.Code)
			}
			body := rec.Body.String()
			if !strings.Contains(body, tc.wantVer) {
				t.Errorf("body no contiene veredicto %q", tc.wantVer)
			}
			if !strings.Contains(body, tc.wantInBody) {
				t.Errorf("body no contiene %q", tc.wantInBody)
			}
		})
	}
}

// TestSalud valida que el health-check responde JSON bien formado con ALLOW.
func TestSalud(t *testing.T) {
	c, rec := newCtx("/midos/salud")
	if err := salud(c); err != nil {
		t.Fatalf("salud() error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	if got["verdict"] != "ALLOW" {
		t.Errorf("verdict: got %v, want ALLOW", got["verdict"])
	}
	if got["zona"] != "midos" {
		t.Errorf("zona: got %v, want midos", got["zona"])
	}
}

// TestInterfaz garantiza en tiempo de compilación + test que Module satisface
// la firma esperada por main.go (Endpoint + Register vía New()).
func TestInterfaz(t *testing.T) {
	m := New()
	if m.Endpoint() != "midos" {
		t.Errorf("Endpoint: got %q, want midos", m.Endpoint())
	}
	// Register no debe panickear al cablear sobre un grupo real.
	e := echo.New()
	g := e.Group("/midos")
	m.Register(g) // si panickea, el test falla con stack
}

// TestLinksNoRotos es la regresión de un 404 real: con la convención
// g.GET("/", ...), la ruta /midos (sin slash final) no existe — ningún
// template debe linkear ahí. La home de la zona es /midos/.
func TestLinksNoRotos(t *testing.T) {
	pages := []struct {
		name    string
		handler func(c *echo.Context) error
	}{
		{"home", home},
		{"coherencia", coherencia},
	}
	for _, p := range pages {
		t.Run(p.name, func(t *testing.T) {
			c, rec := newCtx("/midos/")
			if err := p.handler(c); err != nil {
				t.Fatalf("%s error: %v", p.name, err)
			}
			if strings.Contains(rec.Body.String(), `href="/midos"`) {
				t.Errorf(`%s linkea a /midos (sin slash) — eso es 404; usar /midos/`, p.name)
			}
		})
	}
}
