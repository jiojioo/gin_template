package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, profile)
}
