// Package handler wires HTTP handlers to services.
package handler

import (
	userhandler "github.com/jiojioo/gin_template/internal/handler/v1/user"
	"github.com/jiojioo/gin_template/internal/service"
)

type Handler struct {
	User *userhandler.Handler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{User: userhandler.NewHandler(services.User)}
}
