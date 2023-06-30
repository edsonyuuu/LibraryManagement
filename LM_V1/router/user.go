package router

import (
	"LibraryManagementV1/LM_V1/logic"
	"LibraryManagementV1/LM_V1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func userRouter(r *gin.Engine) {
	//路由就是 /user/users
	base := r.Group("/user")
	base.Use(userCheck())
	user := base.Group("/users")
	{
		user.GET("", logic.GetUser) //
		//user.PUT("/:id", logic.UpdateUser)
		user.PUT("", logic.UpdateUser) //
		//user.DELETE(":id", logic.DeleteUser)
		user.GET("/:id/records", logic.GetUserRecords)               //
		user.GET("/:id/records/:status", logic.GetUserStatusRecords) //
		//用户自助借书还书
		user.POST("/records/:bookId", logic.BorrowBook) //
		user.PUT("/records/:bookId", logic.ReturnBook)  //
	}

	// 暂不使用
	/*book := base.Group("/books")
	{
		book.GET("/:id", logic.GetBook)
		//book.POST("/:id", logic.AddBook)
		//book.DELETE("/:id", logic.DeleteBook)
	}
	category := base.Group("/categories")
	{
		//category.GET("/:id", logic.GetCategory)
		category.GET("/:id/books/:type", logic.GetCategoryBooks)
	}*/
}

func userCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("debug") != "" {
			c.Next()
			return
		}

		auth := c.GetHeader("Authorization")
		fmt.Printf("auth:%+v\n", auth)
		data, err := tools.Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "验签失败",
			})
			return
		}

		fmt.Printf("data:%+v\n", data)
		if data.ID <= 0 || data.Name == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "用户信息获取失败",
			})
			return
		}
		c.Next()
	}
}
