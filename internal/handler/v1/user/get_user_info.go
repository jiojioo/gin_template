package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/service"
	"github.com/jiojioo/gin_template/pkg/response"
)

func (h *Handler) GetUserInfo(c *gin.Context) {
	rawUserID, ok := c.Get("user_id")
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "missing user identity")
		return
	}
	userID, ok := rawUserID.(uint64)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "invalid user identity")
		return
	}
	profile, err := h.users.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(c, http.StatusNotFound, "user not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "unable to load user")
		return
	}
	response.Success(c, profile)
}
