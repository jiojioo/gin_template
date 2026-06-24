package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/service"
	"github.com/jiojioo/gin_template/pkg/response"
)

func (h *Handler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid login request")
		return
	}
	resp, err := h.users.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Fail(c, http.StatusUnauthorized, "invalid username or password")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "login failed")
		return
	}
	response.Success(c, resp)
}
