package logic

import (
	"LibraryManagementV1/LM_V1/model"
	"LibraryManagementV1/LM_V1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type Librarian struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserLogin godoc
//
//	@Summary		用户登录
//	@Description	会执行用户登录操作
//	@Tags			User
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_name	formData	string	true	"用户名"
//	@Param			password	formData	string	true	"密码"
//	@Param			code		formData	string	true	"验证码"
//	@response		200,500		{object}	tools.HttpCode{data=Token}
//	@Router			/userLogin [POST]
func UserLogin(c *gin.Context) {

	code, err := model.RedisConn.Get(c, "code").Result()
	if err != nil {
		fmt.Println("Failed to store verification code in Redis:", err.Error())
	}
	fmt.Println("Redis code:", code)

	// 获取用户输入的验证码
	userCode := c.PostForm("code")
	//验证验证码
	fmt.Println("Redis code:", code)
	fmt.Println("form code:", userCode)
	if userCode == code {

		var user User
		err = c.ShouldBind(&user) //参数绑定，可以绑定json，query，param，yaml，xml，校验不通过，返回错误
		if err != nil {           //[这里有点小疑问？]
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.UserError,
				Message: "您的信息有误",
			})
			fmt.Printf("data:%+v\n", user)
			return
		}

		DbUser := model.GetUser(user.UserName, user.Password)
		fmt.Printf("打印user登录信息:%+v\n", DbUser)

		if DbUser.Id <= 0 {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.NotFound,
				Message: "没有找到用户信息",
				Data:    struct{}{},
			})
			return
		}
		c.SetCookie("id", strconv.FormatInt(DbUser.Id, 10), 3600, "/", "", false, false)
		//下发Token
		a, r, err := tools.Token.GetToken(DbUser.Id, DbUser.UserName)
		fmt.Printf("atoken:%s\n", a)
		fmt.Printf("rtoken:%s\n", r)
		if err != nil {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.UserError,
				Message: "Token生效失败！错误信息：" + err.Error(),
				Data:    struct{}{},
			})
			return
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "登录成功，正在跳转~",
			Data: Token{
				AccessToken:  a,
				RefreshToken: r,
			},
		})
		return
	}

	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.Failed,
		Message: "验证码错误！",
		Data:    struct{}{},
	})
	return
}

// LibrarianLogin godoc
//
//	@Summary		管理员登录
//	@Description	执行管理员登录操作
//	@Tags			Admin登录模块
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			user_name	formData	string	true	"用户名"
//	@Param			password	formData	string	true	"密码"
//	@Param			code		formData	string	true	"验证码"
//	@response		200,500		{object}	tools.HttpCode{data=Token}
//	@Router			/adminLogin [POST]
func LibrarianLogin(c *gin.Context) {
	code, err := model.RedisConn.Get(c, "code").Result()
	if err != nil {
		fmt.Println("Failed to store verification code in Redis:", err.Error())
	}
	fmt.Println("Redis code:", code)

	// 获取用户输入的验证码
	userCode := c.PostForm("code")
	//验证验证码
	fmt.Println("Redis code:", code)
	fmt.Println("form code:", userCode)
	if userCode == code {

		var admin Librarian
		if c.ShouldBind(&admin) != nil {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.Failed,
				Message: "用户信息错误",
			})
			return
		}
		//TODO: 入参校验 和 SQL注入问题
		fmt.Printf("data:%+v\n", admin)
		dbUser := model.CheckAdminMsg(admin.UserName, admin.Password)
		fmt.Printf("user:%+v\n", dbUser)
		if dbUser.Id > 0 {
			// 设置保存session
			err := model.SaveSession(c, dbUser.UserName, dbUser.Id)
			if err != nil {
				c.JSON(http.StatusOK, tools.HttpCode{
					Code:    tools.UserError,
					Message: err.Error(),
				})

			}

			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: "登录成功，正在跳转~",
			})
			return
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.UserError,
			Message: "用户信息错误",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.Failed,
		Message: "验证码错误！",
		Data:    struct{}{},
	})
	return
}

// Logout godoc
//
//	@Summary		用户退出
//	@Description	会执行用户退出操作
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@response		500,302	{object}	tools.HttpCode
//	@Router			/logout [post]
func Logout(c *gin.Context) {
	_ = model.DeleteSession(c)
	c.SetCookie("id", "", 0, "/", "", false, false)
	c.JSON(http.StatusFound, tools.HttpCode{
		Code:    tools.OK,
		Message: "退出登录成功！",
		Data:    struct{}{},
	})
	return
}

// AdminLogout godoc
//
//	@Summary		管理员退出登录
//	@Description	管理员退出登录
//	@Tags			Admin登录模块
//	@Accept			json
//	@Produce		json
//	@response		500,302	{object}	tools.HttpCode
//	@Router			/adminLogout [get]
func AdminLogout(c *gin.Context) {
	_ = model.DeleteSession(c)
	c.JSON(http.StatusFound, tools.HttpCode{
		Code:    tools.OK,
		Message: "退出登录成功！",
	})
	return
}

// SendNum godoc
//
//	@Summary		发送验证码
//	@Description	生成并发送验证码到用户，并将验证码存储在Redis中
//	@Produce		json
//	@Success		200	{object}	tools.HttpCode
//	@Failure		404	{object}	tools.HttpCode
//	@Router			/GetCode [get]
func SendNum(c *gin.Context) {
	// 生成验证码
	SendCode := model.GenerateCode()
	ret := model.RedisConn.Set(c, "code", SendCode, 5*time.Minute)
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "发送验证码成功！",
		})
		fmt.Println("code:", SendCode)
		return
	}
	c.JSON(http.StatusNotFound, tools.HttpCode{
		Code:    tools.Failed,
		Message: "发送验证码失败！",
	})
	return

}

// 原使用session进行登录校验的代码
/*func Login(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user) //参数绑定，可以绑定json，query，param，yaml，xml，校验不通过，返回错误
	if err != nil {            //[这里有点小疑问？]
		c.JSON(200, tools.HttpCode{
			Code:    tools.UserError,
			Message: "您的信息有误",
		})
		fmt.Printf("data:%+v\n", user)
		return
	}

	DbUser := models.GetUser(user.Name, user.Tel, user.Password)
	fmt.Printf("打印user登录信息:%+v\n", DbUser)

	if DbUser.ID > 0 {
		_ = models.SaveSession(c, DbUser.UserName, DbUser.ID)
		if err != nil {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: err.Error(),
			})
		}
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "登录成功，正在跳转~",
		})

		return
	}

	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.UserError,
		Message: "用户信息错误！",
	})

}*/
