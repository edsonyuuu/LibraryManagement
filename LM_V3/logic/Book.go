package logic

import (
	"LibraryManagementV1/LM_V3/model"
	"LibraryManagementV1/LM_V3/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SearchBook godoc
//
//	@Summary		搜索图书
//	@Description	游客获取所有图书信息
//	@Tags			游客模式
//	@Accept			json
//	@Produce		json
//
//	@Param			id			query		string	false	"图书id"
//	@Param			size		query		string	false	"每页书籍数量"
//	@Param			direction	query		string	false	"输入任何内容将向前翻页"
//
//	@response		200,500		{object}	tools.HttpCode{}
//	@Router			/books/page [get]
func SearchBook(c *gin.Context) {

	idStr := c.DefaultQuery("id", "1")
	sizeStr := c.DefaultQuery("size", "100")
	size, err := strconv.Atoi(sizeStr)
	direction := c.DefaultQuery("direction", "back")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "size错误",
		})
		return
	}

	book := model.GetRedisBookId(c, idStr, size, direction)

	if book == nil {
		c.JSON(http.StatusNotFound, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "未查到数据",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "成功！",
		Data:    book,
	})
	return
}

//	@Param	currentPage	query	string	true	"当前页"
//	@Param	pageSize	query	string	true	"页大小"

// SearchCategory godoc
//
//	@Summary		搜索图书种类
//	@Description	获取所有图书种类信息
//	@Tags			游客模式
//	@Accept			json
//	@Produce		json
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/categories [get]
func SearchCategory(c *gin.Context) {
	category := model.GetCategory()
	if category != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code: tools.OK,
			Data: category,
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "未找到图书种类信息！",
		Data:    struct{}{},
	})
	return
}

// GetBook godoc
//
//	@Summary		获取图书bookId的信息
//	@Description	根据图书ID获取图书信息
//	@Tags			游客模式
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int	true	"图书ID"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/books/{id} [get]
func GetBook(c *gin.Context) {
	bookIdStr := c.Param("id")
	fmt.Printf("bookIdStr:%s\n", bookIdStr)
	bookId, _ := strconv.ParseInt(bookIdStr, 10, 64)

	ret := model.GetBook(bookId)
	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "查询图书信息失败！",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "您所查找到的图书信息！",
		Data:    ret,
	})
	return
}

//func FromBookName(c *gin.Context) {
//	bookNameStr := c.Param("name")
//	fmt.Printf("bookIdname:%s\n", bookNameStr)
//	//
//	/*var bookId int64
//	sql := "select id from book where name = ?"
//	err := model.DB.Raw(sql, bookNameStr).Scan(&bookId).Error
//	if err != nil {
//		fmt.Printf("select bookId failed :%+v\n", err.Error())
//	}
//	bookIdStr := strconv.FormatInt(bookId, 10)*/
//	//
//	//
//	ret := model.GetRedisBookId(c, bookNameStr)
//	if ret == nil {
//		c.JSON(http.StatusOK, tools.HttpCode{
//			Code:    tools.NotFound,
//			Message: "查询图书信息失败！",
//			Data:    struct{}{},
//		})
//		return
//	}
//	c.JSON(http.StatusOK, tools.HttpCode{
//		Code:    tools.OK,
//		Message: "您所查找到的图书信息！",
//		Data:    ret,
//	})
//	return
//}

// GetRecords godoc
//
//	@Summary		查询记录表信息
//	@Description	查询记录表中的所有记录
//	@Tags			AdminGet
//	@Produce		json
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/records [get]
func GetRecords(c *gin.Context) {
	ret := model.AdminGetRecords()
	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "record表记录未被找到！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "管理员查询记录表成功！",
		Data:    ret,
	})
	return
}

// GetUserRecordStatus godoc
//
//	@Summary		查询记录表归还或未归还状态信息
//	@Description	根据归还或未归还状态查询记录表信息
//	@Tags			AdminGet
//	@Produce		json
//	@Param			status	path		int	true	"归还或未归还状态，1为未归还，0为已归还"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/records/{status} [get]
func GetUserRecordStatus(c *gin.Context) {
	statusStr := c.Param("status")
	status, _ := strconv.Atoi(statusStr)
	ret := model.GetUserRecordStatus(status)
	if ret == nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "record表归还或未归还状态查询失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "管理员查询记录表归还或未归还状态成功！",
		Data:    ret,
	})
	return
}

// AddBook godoc
//
//	@Summary		添加图书
//	@Description	添加新图书
//	@Tags			Admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			categoryId	formData	int		true	"图书分类ID"
//	@Param			bn			formData	string	true	"图书编号"
//	@Param			name		formData	string	true	"图书名称"
//	@Param			description	formData	string	false	"图书描述"
//	@Param			count		formData	int		true	"图书数量"
//	@response		200,500		{object}	tools.HttpCode
//	@Router			/admin/books [post]
func AddBook(c *gin.Context) {
	CategoryIdStr := c.PostForm("categoryId")
	CategoryId, _ := strconv.ParseInt(CategoryIdStr, 10, 64)
	//
	bn := c.PostForm("bn")
	name := c.PostForm("name")
	description := c.PostForm("description")
	//
	countStr := c.PostForm("count")
	count, _ := strconv.Atoi(countStr)
	ret := model.AddBooks(CategoryId, bn, name, description, count)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "添加图书失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "添加图书成功！",
		Data:    ret,
	})
	return
}

// UpdateBook godoc
//
//	@Summary		更新图书信息
//	@Description	根据图书ID更新图书信息
//	@Tags			AdminPUT
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		int		true	"图书ID"
//	@Param			bn			formData	string	false	"图书编号"
//	@Param			name		formData	string	false	"图书名称"
//	@Param			description	formData	string	false	"图书描述"
//	@Param			count		formData	int		false	"图书数量"
//	@Param			categoryId	formData	int		false	"图书种类ID"
//	@Success		200,500		{object}	tools.HttpCode
//	@Router			/admin/books/{id} [put]
func UpdateBook(c *gin.Context) {
	// 根据图书id来更新图书信息
	bookIdStr := c.Param("id")
	bookId, _ := strconv.ParseInt(bookIdStr, 10, 64)
	//
	bn := c.PostForm("bn")
	name := c.PostForm("name")
	description := c.PostForm("description")
	countStr := c.PostForm("count")
	count, _ := strconv.Atoi(countStr)
	//
	categoryIdStr := c.PostForm("categoryId")
	categoryId, _ := strconv.ParseInt(categoryIdStr, 10, 64)

	ret := model.UpdateBooks(bookId, bn, name, description, count, categoryId)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "更新失败",
		})
		fmt.Printf("管理员更新图书信息失败！err:%d\n", bookId)
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "图书信息更新成功！",
	})
	return
}

// DeleteBook godoc
//
//	@Summary		删除图书
//	@Description	根据图书ID删除图书信息
//	@Tags			AdminDelete
//	@Produce		json
//	@Param			id		path		int	true	"图书ID"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	bookIdStr := c.Param("id")
	bookId, _ := strconv.ParseInt(bookIdStr, 10, 64)
	ret := model.DeleteBook(bookId)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "删除图书失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "该图书已被删除！",
	})
	return
}

// GetCategory 管理员获取图书种类
func GetCategory(c *gin.Context) {

}

// AddCategory godoc
//
//	@Summary		添加图书种类信息
//	@Description	添加新的图书种类信息
//	@Tags			Admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData	string	true	"图书种类名称"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/categories [post]
func AddCategory(c *gin.Context) {
	name := c.PostForm("name")
	ret := model.AddCategory(name)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "添加图书种类信息失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "添加图书种类信息成功！",
		Data:    struct{}{},
	})
	return
}

// UpdateCategory
//
//	@Summary		更新图书种类信息
//	@Description	根据图书分类ID更新图书分类信息
//	@Tags			AdminPUT
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int		true	"图书分类ID"
//	@Param			name	formData	string	true	"分类名称"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	categoryIdStr := c.Param("id")
	categoryId, _ := strconv.ParseInt(categoryIdStr, 10, 64)
	name := c.PostForm("name")

	ret := model.UpdateCategory(categoryId, name)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "更新图书种类信息失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "更新图书种类信息成功！",
		Data:    struct{}{},
	})
	return

}

// DeleteCategory godoc
//
//	@Summary		删除图书种类
//	@Description	根据图书种类ID删除图书种类信息
//	@Tags			AdminDelete
//	@Produce		json
//	@Param			id		path		int	true	"图书种类ID"
//	@response		200,500	{object}	tools.HttpCode
//	@Router			/admin/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	categoryIdStr := c.Param("id")
	categoryId, _ := strconv.ParseInt(categoryIdStr, 10, 64)
	ret := model.DeleteCategory(categoryId)
	if ret == false {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.Failed,
			Message: "删除图书种类失败！",
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "该图书种类已被删除！",
	})
	return
}
