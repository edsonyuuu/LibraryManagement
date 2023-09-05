package logic

import (
	"LibraryManagementV1/LM_V3/global"
	"LibraryManagementV1/LM_V3/model"
	"LibraryManagementV1/LM_V3/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
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
//	@Param			phone		formData	string	true	"手机号"
//	@Param			code		formData	string	true	"验证码"
//	@response		200,500		{object}	tools.HttpCode{data=Token}
//	@Router			/userLogin [POST]
func UserLogin(c *gin.Context) {
	phone := c.PostForm("phone")
	ans := IsPhoneNumber(phone)
	if ans == false {
		ans := IsPhoneNumber(phone)
		if ans == false {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.Failed,
				Message: "输入手机号为非法手机号，请重新输入~",
			})
			c.Abort()
			return
		}
	}
	code, err := global.RedisConn.Get(c, phone).Result()
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
//	@Param			phone		formData	string	true	"手机号"
//	@Param			code		formData	string	true	"验证码"
//	@response		200,500		{object}	tools.HttpCode{data=Token}
//	@Router			/adminLogin [POST]
func LibrarianLogin(c *gin.Context) {
	phone := c.PostForm("phone")
	ans := IsPhoneNumber(phone)
	if ans == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "输入手机号为非法手机号，请重新输入~",
		})
		c.Abort()
		return
	}
	code, err := global.RedisConn.Get(c, phone).Result()
	if err != nil {
		fmt.Println("Failed to store verification code in Redis:", err.Error())
	}
	fmt.Println("Redis code:", code)
	// 获取用户输入的验证码
	userCode := c.PostForm("code")
	//验证验证码

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
	c.JSON(http.StatusOK, tools.HttpCode{
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
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/adminLogout [get]
func AdminLogout(c *gin.Context) {
	_ = model.DeleteSession(c)
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "退出登录成功！",
	})
	return
}

// SendNum 此方法已被启用
func SendNum(c *gin.Context) {
	// 生成验证码
	SendCode := model.GenerateCode()
	ret := global.RedisConn.Set(c, "code", SendCode, 5*time.Minute)
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

type verificationCode struct {
	Phone string
	Count int
}

var verificationCodes map[string]*verificationCode

// AliSendMsg godoc
//
//	@Summary		发送验证码
//	@Description	生成并发送验证码到用户，并将验证码存储在Redis中
//	@Produce		json
//	@Param			phone	path		string	true	"手机号"
//	@Success		200		{object}	tools.HttpCode
//	@Failure		404		{object}	tools.HttpCode
//	@Router			/GetCode/{phone} [get]
func AliSendMsg(c *gin.Context) {
	phone := c.Param("phone")
	// 判断手机号是否是非法手机号
	ans := IsPhoneNumber(phone)
	if ans == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "输入手机号为非法手机号，请重新输入~",
		})

		return
	}
	/*if verificationCodes[phone].Count >= 3 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "发送到此手机号的验证码的次数今日已达上限~",
		})
		return
	}*/
	// 这里执行的是从redis中获取5分钟内未过期的短信验证码，若上边次数限制的代码没有被注释掉，可以使用下方代码，从redis中拿验证码
	// 可以减少发送验证码的次数，只是为测试使用而减少
	/*msg, err := model.RedisConn2.Get(c, phone).Result()
	if err != nil {
		fmt.Printf("从redis中获取未过期的短信验证码出错~ err:%+v\n", err.Error())
	}
	if msg != "" {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "短信验证码还未过期",
			//Data:    msg,
		})
		c.Abort()
		return
	}*/
	//
	// 随机生成六位数验证码
	SendCode := model.GenerateCode()
	// 调用阿里云API执行发送验证码到手机的功能
	//tools.Aliyun(phone, SendCode)
	// 设置验证码在redis中的缓存时间,5分钟，3*24*60
	ret := global.RedisConn.Set(c, phone, SendCode, 5*time.Minute)

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

func IsPhoneNumber(input string) bool {
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(input)
}
