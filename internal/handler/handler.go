package handler

import (
	v1 "Iris/internal/handler/v1"
	"Iris/internal/service"
	"net/http"
	"text/template"

	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

// NewHandler creates and returns an http.Handler with API routes,
// static file serving and the web frontend at the root path.
func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())
	handler.Static("/static", "./web/static")

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/shorten", handlerV1.Shorten)

	apiV1.GET("/s/:short_url", handlerV1.Redirect)
	apiV1.GET("/analytics/:short_url", handlerV1.GetAnalytics)

	handler.GET("/", homePage(template.Must(template.ParseFiles(templatePath))))

	return handler

}

// homePage renders the main HTML page.
func homePage(tmpl *template.Template) func(c *ginext.Context) {
	return func(c *ginext.Context) {
		c.Header("Content-Type", "text/html")
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
		}
	}
}
