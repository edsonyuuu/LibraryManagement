package router

import (
	_ "LibraryManagementV1/LM_V4/docs"
	"LibraryManagementV1/LM_V4/global"
	"LibraryManagementV1/LM_V4/logic"
	"LibraryManagementV1/LM_V4/model"
	"LibraryManagementV1/LM_V4/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func New() *gin.Engine {
	// cd .\LM_V2\
	// http://localhost:8083/swagger/index.html
	r := gin.Default()
	//协程执行定时器
	//go DingShi()
	//go PreDingShi()
	//
	//r.Use(limitClick(5, 10))
	userRouter(r)
	adminRouter(r)
	r.Static("/kit", "./kit")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//验证码
	r.GET("/GetCode/:phone", logic.AliSendMsg)
	r.POST("/userLogin", logic.UserLogin)
	r.POST("/userLogout", logic.Logout)
	r.POST("/users", logic.AddUser)
	r.POST("/adminLogin", logic.LibrarianLogin)
	r.GET("/adminLogout", logic.AdminLogout)
	//游客可以浏览书籍和分类
	r.GET("/books/page", logic.SearchBook)
	//r.GET("/books/:id", logic.GetBook)
	r.GET("/books/:id", logic.GetBook)
	//r.GET("/book/BookName")
	r.GET("/categories", logic.SearchCategory)

	return r
}

func DingShi() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	//
	for {
		select {
		//定时器触发操作
		case <-t.C:
			model.Listen()
		}
	}
}

// PreDingShi 预热数据，将前几页数据添加进redis中的定时器
func PreDingShi() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	id := int64(1)
	size := 100
	for {
		select {
		case <-t.C:
			books := model.PreHeating(id, size)
			model.SavePreHeatingBooks(id, books)
			id += int64(size)
			if id == int64(301) {
				id = int64(1)
			}
		}
	}
}

// IP+UA middleware
func limitClick(maxCount int, t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取IP
		ip := c.ClientIP()
		//获取UA
		ua := c.GetHeader("User-Agent")
		url := strings.Split(c.Request.URL.Path+"?"+c.Request.URL.RawQuery, "/")
		lastUrl := url[len(url)-1:]
		pathStr := fmt.Sprintf("%v_%v_%v", ip, ua, lastUrl)
		fmt.Println(pathStr)
		//
		requestQuery := global.RedisConn
		//访问的路径次数加一
		requestQuery.Incr(c, pathStr)
		//
		requestCountStr, _ := requestQuery.Get(c, pathStr).Result()
		requestCount, _ := strconv.Atoi(requestCountStr)
		//
		if requestCount > maxCount {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.Failed,
				Message: "请求太频繁，请稍候再试~",
				Data:    nil,
			})
			c.Abort()
			return
		} else if requestCount == 1 {
			err := requestQuery.Expire(c, pathStr, t*time.Second).Err()
			if err != nil {
				fmt.Printf("设置过期时间错误err:%+v\n", err.Error())
				c.JSON(http.StatusOK, tools.HttpCode{
					Code:    tools.Failed,
					Message: "设置过期时间错误",
				})
				c.Abort()
				return
			}

		}
		c.Next()
	}
}

func CheckRecord() {
	//t := time.NewTicker(10 * time.Minute)
	now := time.Now()
	fmt.Printf("现在的时间点：%+v\n", now)
	next := now.Add(5 * time.Minute)
	fmt.Printf("打印next:%+v\n", next)

	//
	ret := model.ReturnTime(now)
	fmt.Println(ret)
	//
	ans := model.BanUsers(ret)
	fmt.Printf("打印ans%+v\n", ans)
}

func WillReturn() {
	now := time.Now()

	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	// 计算当前时间与下一个零点的时间差
	duration := next.Sub(now)
	time.Sleep(duration)
	//
	model.AdvanceOneDay()

}
