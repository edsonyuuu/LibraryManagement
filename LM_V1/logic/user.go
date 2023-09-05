package logic

import (
	"LibraryManagementV1/LM_V1/model"
	"LibraryManagementV1/LM_V1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetUser godoc
//
//	@Summary		获取用户信息
//	@Description	根据用户ID获取用户信息
//	@Tags			用户信息
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users [get]
func GetUser(c *gin.Context) {
	userIdStr, _ := c.Cookie("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	ret := model.LookSelfMsg(userId)
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "查看用户信息成功！",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "ret中未查到信息",
		Data:    struct{}{},
	})
	return
}

// UpdateUser 更新用户信息
//
//	@Summary		更新用户信息
//	@Description	根据请求中的参数更新用户信息
//	@Tags			用户信息
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@Param			user_name		formData	string	false	"新用户名"
//	@Param			password		formData	string	true	"旧密码"
//	@Param			new_password	formData	string	true	"新密码"
//	@Param			phone			formData	string	false	"新手机号"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users [put]
func UpdateUser(c *gin.Context) {
	var user model.User

	err := c.ShouldBind(&user) //参数绑定，可以绑定json，query，param，yaml，xml，校验不通过，返回错误
	if err != nil {            //[这里有点小疑问？]
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.UserError,
			Message: "您的信息有误",
		})
		fmt.Printf("UpdateUserData:%+v\n", user)
		return
	}

	userIdStr, _ := c.Cookie("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	userName := c.PostForm("user_name")       //新用户名
	password := c.PostForm("password")        //旧密码
	newPassword := c.PostForm("new_password") //新密码
	phone := c.PostForm("phone")              //新手机号

	if user.Password == password {
		ret := model.UpdateUserMsg(userId, userName, newPassword, phone)

		if ret == nil {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.Failed,
				Message: "用户更新信息失败！",
				Data:    struct{}{},
			})
			return
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "用户信息更新成功！",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.Failed,
		Message: "原密码错误，无法更新用户信息！",
		Data:    struct{}{},
	})
	return
}

// GetUserRecords godoc
//
//	@Summary		获取用户借阅记录
//	@Description	根据用户ID查询其借阅记录
//	@Tags			用户借阅记录
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		int		true	"用户ID"
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users/{id}/records [get]
func GetUserRecords(c *gin.Context) {
	//去看用户id为？的借阅记录
	userIdStr := c.Param("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	ret := model.UserGetRecords(userId)

	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "ret为空",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "查询用户借阅记录成功",
		Data:    ret,
	})
	return
}

// GetUserStatusRecords godoc
//
//	@Summary		获取用户借阅记录状态
//	@Description	根据用户ID和status来查看某个用户没有还书或还书的情况
//	@Tags			用户借阅或归还状态
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		int		true	"用户ID"
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users/{id}/records/{status} [get]
func GetUserStatusRecords(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	statusStr := c.Param("status")
	//
	status, _ := strconv.Atoi(statusStr)
	ret := model.GetUserBorrowStatus(userId, status)
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "查询用户借阅状态成功",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "查询用户借阅状态为空",
	})
	return
}

// BorrowBook godoc
//
//	@Summary		借阅图书
//	@Description	根据用户ID和图书ID进行借阅操作
//	@Tags			图书借阅与归还
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			bookId			path		int		true	"图书ID"
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users/records/{bookId} [post]
func BorrowBook(c *gin.Context) {
	userIdStr, _ := c.Cookie("id")
	bookIdStr := c.Param("bookId")
	fmt.Printf("bookIdStr:%s\n", bookIdStr)
	fmt.Printf("userIdStr:%s\n", userIdStr)
	//
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	borrowId, _ := strconv.ParseInt(bookIdStr, 10, 64)
	//

	ret := model.CreateRecord(userId, borrowId)
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "借阅成功！",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "ret记录未找到",
		Data:    struct{}{},
	})
	return
}

// ReturnBook godoc
//
//	@Summary		还书
//	@Description	用户归还图书
//	@Tags			图书借阅与归还
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			bookId			path		int		true	"图书ID"
//	@Param			Authorization	header		string	false	"Bearer 用户令牌"
//	@response		200,500			{object}	tools.HttpCode
//	@Router			/user/users/records/{bookId} [put]
func ReturnBook(c *gin.Context) {
	userIdStr, _ := c.Cookie("id")
	bookIdStr := c.Param("bookId")

	fmt.Printf("bookIdStr:%s\n", bookIdStr)
	fmt.Printf("userIdStr:%s\n", userIdStr)
	//
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	bookId, _ := strconv.ParseInt(bookIdStr, 10, 64)

	//
	model.ReturnBook(userId, bookId)

	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "还书成功！",
	})
	return

}
