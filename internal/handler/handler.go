package handler

import (
	v1 "Iris/internal/handler/v1"
	"Iris/internal/service"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/shorten", handlerV1.Shorten)

	apiV1.GET("/s/:short_url", handlerV1.Redirect)
	apiV1.GET("/analytics/:short_url", handlerV1.GetAnalytics)

	return handler

}
