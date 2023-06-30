package logic

import (
	"LibraryManagementV1/LM_V2/model"
	"LibraryManagementV1/LM_V2/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SearchUser godoc
//
// @Summary 查询所有用户
// @Description 查看所有用户信息
// @Tags  AdminGet
// @Produce json
// @Success 200,500 {object} tools.HttpCode
// @Router /admin/users [get]
func SearchUser(c *gin.Context) {
	idStr, _ := c.Cookie("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotLogin,
			Message: "管理员未登录！",
		})
		return
	}
	// 查看所有用户
	ret := model.GetUsers()
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "管理员查询所有用户信息成功！",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "管理员查询所有用户信息失败！",
	})
	return
}

// AddUser godoc
//
// @Summary 用户注册
// @Description 用户注册
// @Tags  User
// Accept  multipart/form-data
// @Produce json
// @Param user_name formData string true "用户名"
// @Param phone formData string true "手机号码"
// @Param password formData string true "用户密码"
// @Param name formData string true "用户真实姓名"
// @Param sex formData string true "性别"
// @response 200,500 {object} tools.HttpCode
// @Router /users [post]
func AddUser(c *gin.Context) {
	//获取参数
	userName := c.PostForm("user_name") //注册用户名
	fmt.Printf("userName:%+v", userName)

	userTel := c.PostForm("phone") //注册的手机号
	fmt.Printf("userTel:%+v", userTel)

	passWord := c.PostForm("password") //注册用户名密码
	fmt.Printf("passWord:%+v\n", passWord)

	name := c.PostForm("name")
	fmt.Printf("name:%+v\n", name)

	sex := c.PostForm("sex")
	fmt.Printf("sex:%+v\n", sex)

	//数据验证
	if len(userName) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.NotRegister,
			"msg":  "用户名不能为空！",
		})
		return
	}
	if len(userTel) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.NotRegister,
			"msg":  "注册手机号不能为空！",
		})
		return
	}
	if len(passWord) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.NotRegister,
			"msg":  "密码不能设置为空！",
		})
		return
	}
	if len(name) <= 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "name不能为空！",
		})
		return
	}

	if sex == "男" && sex == "女" {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "错误的性别",
		})
		return
	}

	//判断手机号是否存在
	var user model.User
	sql1 := "SELECT * FROM user WHERE phone = ? LIMIT 1"
	err := model.DB.Raw(sql1, userTel).Scan(&user).Error

	if user.Id != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": tools.Failed,
			"msg":  "用户名已存在！",
		})
		return
	}

	sql2 := fmt.Sprintf("INSERT INTO user (user_name, password, phone,name,sex) VALUES ('%s', '%s', '%s','%s','%s')", userName, passWord, userTel, name, sex)
	err = model.DB.Exec(sql2).Error

	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": tools.OK,
		"msg":  "用户注册成功！",
	})
	return
}

// UpdateUserByAdmin 管理员修改用户信息
// @Summary 管理员修改用户信息
// @Description 管理员修改用户信息
// @Tags AdminPUT
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "UserID"
// @Param user_name formData string true "更新用户名"
// @Param newPwd formData string true "更新用户密码"
// @Param phone formData string true "更新用户手机号码"
// @response 200,500 {object} tools.HttpCode
// @Router /admin/users/{id} [put]
func UpdateUserByAdmin(c *gin.Context) {
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

	userIdStr := c.Param("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	//
	userName := c.PostForm("user_name")
	//password := c.PostForm("password")
	newPwd := c.PostForm("newPwd")
	phone := c.PostForm("phone")

	// 使用用户更新身份信息的方法
	ret := model.UpdateUserMsg(userId, userName, newPwd, phone)
	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "管理员更新用户信息失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "管理员更新用户信息成功！",
		Data:    ret,
	})
	return
	//if password == user.Password {

	//}
	//c.JSON(http.StatusOK, tools.HttpCode{
	//	Code:    tools.Failed,
	//	Message: "用户原密码错误",
	//})
	//return

}

// DeleteUser godoc
//
// @Summary 删除用户信息
// @Description 根据用户ID删除用户信息
// @Tags AdminDelete
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "用户ID"
// @Success 200,500 {object} tools.HttpCode
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	ret := model.DeleteUser(userId)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "删除用户信息失败！",
			Data:    ret,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "删除用户信息成功！",
		Data:    ret,
	})
	return
}

// GetUserBook 获取用户已归还或者未归还的所有记录
// @Summary 获取用户借阅状态
// @Description 管理员获取用户借阅状态
// @Tags AdminGet
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "用户ID"
// @Param status path int true "借阅状态码（1为借阅；0为未借阅）"
// @response 200,500 {object} tools.HttpCode
// @Router /admin/users/{id}/records/{status} [get]
func GetUserBook(c *gin.Context) {
	//    /admin/users/:id/records/:status
	userIdStr := c.Param("id")
	statusStr := c.Param("status")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	status, _ := strconv.Atoi(statusStr)
	//
	fmt.Printf("userIdStr:%+v\n", userIdStr)
	fmt.Printf("statusStr:%+v\n", statusStr)
	//
	ret := model.GetUserBorrowStatus(userId, status)
	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查取用户归还状态失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "查取用户借阅状态成功！",
		Data:    ret,
	})
	return
}
