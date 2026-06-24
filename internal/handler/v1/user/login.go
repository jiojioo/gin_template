package user

import (
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
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(c, resp)
}
