package handler

import (
	"Iris/internal/service"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	// apiV1 := handler.Group("/api/v1")
	// handlerV1 := v1.NewHandler(service)

	return handler

}
