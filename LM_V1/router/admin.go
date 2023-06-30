package router

import (
	"LibraryManagementV1/LM_V1/logic"
	"LibraryManagementV1/LM_V1/model"
	"LibraryManagementV1/LM_V1/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

func adminRouter(r *gin.Engine) {
	//librarian := r.Group("/librarians").Use(librarianCheck())
	//      /admin/users
	base := r.Group("/admin")
	base.Use(librarianCheck())
	user := base.Group("/users")
	{
		//user.GET("/:id", logic.GetUser)
		user.GET("", logic.SearchUser)            //
		user.PUT("/:id", logic.UpdateUserByAdmin) //
		user.DELETE("/:id", logic.DeleteUser)     //
		//获取该用户已归还或者未归还的所有记录
		user.GET("/:id/records/:status", logic.GetUserBook) //
		//user.POST("/:id/records/:bookId", logic.BorrowBook)
		//user.PUT("/:id/records/:bookId", logic.ReturnBook)
	}
	//书的所有资源
	//    /admin/books
	book := base.Group("/books")
	{
		book.GET("/:id", logic.GetBook) // 直接使用谁都可以查看图书，此路径先不用
		//book.GET("", logic.SearchBook)
		book.POST("", logic.AddBook)          //
		book.PUT("/:id", logic.UpdateBook)    //
		book.DELETE("/:id", logic.DeleteBook) //
	}

	//   /admin/categories
	category := base.Group("/categories")
	{
		category.GET("/:id", logic.GetCategory) //这个不必要写
		//category.GET("", logic.SearchCategory)
		category.POST("", logic.AddCategory)          //
		category.PUT("/:id", logic.UpdateCategory)    //
		category.DELETE("/:id", logic.DeleteCategory) //
	}
	//记录表的资源  /admin/records
	record := base.Group("/records")
	{
		//所有借书还书记录
		record.GET("", logic.GetRecords) //
		//所有归还或者未归还的记录
		record.GET("/:status", logic.GetUserRecordStatus) //
	}
}
func librarianCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := model.GetSession(c)
		id, ok1 := data["id"]
		name, ok2 := data["name"]
		idInt64, _ := id.(int64)
		if !ok1 || !ok2 || idInt64 <= 0 || name == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "管理员信息获取失败",
			})
			return
		}
		//c.Next()

		c.Set("name", name)
		c.Set("id", idInt64)
		c.Next()
	}
}
