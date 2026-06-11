package saburou

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v5"
)

//go:embed index.gohtml
var indexTemplate string

//go:embed assets/styles.css
var stylesCSS []byte

//go:embed assets/*.png assets/*.jpeg
var images embed.FS

var tmpl = template.Must(template.New("index").Parse(indexTemplate))

// Person representa a alguien a quien agradecer en la página.
type Person struct {
	Name     string // Nombre visible, no tiene por qué coincidir con el nombre de la imagen.
	Image    string // Fichero dentro de assets/.
	Message  string
	URL      string // Opcional.
	Endpoint string // Opcional, ruta a su propia página dentro del sitio (ej. "/chavsi").
	Active   bool   // Si aún no tiene carpeta/módulo propio, se deja desactivado.
}

var thanks = []Person{
	{
		Name:     "Chavsi",
		Image:    "chavsi.png",
		Message:  "Gracias por apoyar la idea del proyecto común y aportar tus ideas! Siento mucho haberte sacado de Flask! :3",
		URL:      "https://github.com/Chavsi",
		Endpoint: "/chavsi/",
		Active:   true,
	},
	{
		Name:    "Dixion",
		Image:   "dixion.jpeg",
		Message: "Gracias por el regalo que me hiciste, gracias a el he podido terminar mi sección muy rápido!",
		URL:     "https://github.com/luisgabrielroldan",
		Active:  true,
	},
	{
		Name:     "Kagliostro",
		Image:    "kagliostro.png",
		Message:  "Gracias por esforzarte tanto por hacerme un regalo y por ser mi amigo tantos años ya! Algún espero que nos encontremos en tierras argentinas.",
		URL:      "https://github.com/ginobadhouse",
		Endpoint: "/kagliostro/",
		Active:   true,
	},
	{
		Name:     "Killa",
		Image:    "killa.png",
		Message:  "Gracias por existir y por conseguir nuevamente que mi cumple sea algo inesperadamente grande. Te quiero mucho <3",
		URL:      "https://github.com/1337xxxLucyF3rxxx1337",
		Endpoint: "/ladythekilla/",
		Active:   true,
	},
	{
		Name:     "Leandro",
		Image:    "leandro.png",
		Message:  "Gracias por apoyarme en cada uno de los problemas que he tenido y que te he contado y por dejarme usar MidOS de vez en cuando. ¡Eres una gran persona!",
		URL:      "https://midos.dev/",
		Endpoint: "/midos/",
		Active:   true,
	},
	{
		Name:     "Lein",
		Image:    "lein.png",
		Message:  "Gracias por crear una página tan guay para mi con tantos efectos especiales y por ayudarme a mi y a los demás con tantos problemas que hemos tenido.",
		URL:      "https://github.com/CalumRakk",
		Endpoint: "/leinSeab/",
		Active:   true,
	},
	{
		Name:     "Saburou",
		Image:    "saburou.png",
		Message:  "¡Para que tengamos muchos más eventos como este!",
		URL:      "https://github.com/heizeisaburou",
		Endpoint: "/saburou/",
		Active:   true,
	},
	{
		Name:     "Savage",
		Image:    "savage.png",
		Message:  "Gracias por la tarta y por evitar que pierda el timepo convenciendome de no hacer implementaciones imposibles en go.",
		URL:      "https://github.com/lsproule",
		Endpoint: "/savage/",
		Active:   true,
	},
}

type Module struct{}

func New() Module {
	return Module{}
}

func (Module) Endpoint() string {
	return "saburou"
}

func (Module) Register(g *echo.Group) {
	g.GET("/", Home)
	g.GET("", Home)
	g.GET("/styles.css", serveCSS)
	g.GET("/assets/:file", serveImage)
}

// Home renderiza la página principal de saburou. Se reutiliza también
// como punto de entrada del sitio ("/") en main.go.
func Home(c *echo.Context) error {
	active := make([]Person, 0, len(thanks))
	for _, p := range thanks {
		if p.Active {
			active = append(active, p)
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ Thanks []Person }{Thanks: active}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

func serveCSS(c *echo.Context) error {
	return c.Blob(http.StatusOK, "text/css; charset=utf-8", stylesCSS)
}

func serveImage(c *echo.Context) error {
	name := filepath.Base(c.Param("file"))

	data, err := images.ReadFile("assets/" + name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "imagen no encontrada")
	}

	contentType := "image/png"
	if strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".jpg") {
		contentType = "image/jpeg"
	}

	return c.Blob(http.StatusOK, contentType, data)
}
