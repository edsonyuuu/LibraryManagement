package logic

import (
	"LibraryManagementV1/LM_V3/model"
	"LibraryManagementV1/LM_V3/tools"
	"github.com/gin-gonic/gin"
	"net/http"

	"strconv"
)

// SendMsg godoc
//
// @Summary		用户收件箱
// @Description	获取用户收件箱信息(这个接口获取消息没有什么意义，只是提示有书没还)
// @Tags		user/users
// @Produce		json
// @CookieParam id string true "用户id"
// @Param Authorization header string true "用户令牌"
// @response 200,500 {object} tools.HttpCode
// @Router			/user/users/messages [GET]
func SendMsg(c *gin.Context) {
	userIdStr, _ := c.Cookie("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	msg := model.GetMessage(userId)
	if len(msg) > 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "您有书没还！",
			Data:    msg,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "您的书已还过！",
		Data:    nil,
	})
	return
}
