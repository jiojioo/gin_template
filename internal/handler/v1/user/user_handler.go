package user

import "github.com/jiojioo/gin_template/internal/service"

type Handler struct {
	users *service.UserService
}

func NewHandler(users *service.UserService) *Handler {
	return &Handler{users: users}
}
