package model

import (
	"LibraryManagementV1/LM_V4/global"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GetCategory() []*Category {
	category := make([]*Category, 0)
	sql1 := "SELECT * FROM `category`"
	err := global.DB.Raw(sql1).Scan(&category).Error

	if err != nil {
		fmt.Printf("查询图书种类err:%+v\n", err.Error())
	}
	return category
}

// GetBooks  此函数方法获取前端传来参数  // offset 方法暂时不适用,此方法暂时不用
func GetBooks(c *gin.Context, pageStr string, page int, size int) []Book {
	key := "book" + pageStr //pageStr 为查询页
	fmt.Println(key)
	//根据传入的pageStr构造键名key，用于在Redis中查找数据,并将数据保存在data中.
	data, err := global.RedisConn.Get(c, key).Bytes()
	//如果Redis中不存在该数据，即err为redis.Nil，则说明需要从数据库中获取数据，并进行压缩存储到Redis中。
	if err == redis.Nil {
		book := make([]Book, 0)
		err = global.DB.Offset((page - 1) * size).Limit(size).Find(&book).Error //分页查询
		if err != nil {
			fmt.Printf("err:%+v\n", err.Error())
		}
		//进行压缩
		Y := YS(book)
		data = Y
		//使用RedisConn.Set方法将key和Y存储到Redis中，并设置过期时间为5小时。
		err = global.RedisConn.Set(c, key, Y, 5*time.Minute).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			fmt.Println("压缩错误：Y")
			return nil
		}
	}

	//调用SelectBooks函数对data进行解压和反序列化，将解压后的数据返回。
	//序列化（Serialization）是指将数据结构或对象状态转换为可存储或传输的格式的过程，例如将对象转换为字节序列或字符序列。
	//反序列化（Deserialization）则是将序列化的数据重新转换为原始数据结构或对象的过程。
	return SelectBooks(data)
}

// SelectBooks 解压缩
func SelectBooks(key []byte) []Book {

	// 创建 gzip 解码器,并将传入的key作为输入源与gzipReader关联起来，用于读取压缩后的JSON数据。
	gzipReader, err := gzip.NewReader(bytes.NewReader(key))
	if err != nil {
		fmt.Printf("vgzip.NewReader 时出现错误！err:%+v\n", err.Error())
		return nil
	}
	//延迟关闭gzipReader，最后释放资源
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			fmt.Printf("gzipReader.Close() 时出现错误！err:%+v\n", err.Error())
			return
		}
	}(gzipReader)
	//
	books := make([]Book, 0)
	//
	// 分批读取解压缩后的 JSON 数据
	//使用json.NewDecoder方法创建一个JSON解码器decoder，将其与gzipReader关联起来。
	decoder := json.NewDecoder(gzipReader)
	batchSize := 100 // 每次读取的批次大小为 100 条记录
	//
	for {
		//每次读取一批数据，反序列化为一个Book类型的切片batch，并将其添加到books中。
		var batch []Book
		err := decoder.Decode(&batch)
		if err == io.EOF { // 已经读取完数据
			break
		} else if err != nil {
			fmt.Printf("ecoder.Decode(&batch) 时出现错误！err:%+v\n", err.Error())
			return nil
		}
		books = append(books, batch...)

		if len(books) >= batchSize {
			break
		}
	}
	//对books中的每个元素的ImgUrl字段进行处理，将其修改为"/kit/" + books[i].ImgUrl。
	for i := 0; i < len(books); i++ {
		books[i].ImgUrl = "/kit/" + books[i].ImgUrl
	}
	return books
}

// GetRedisBookId 将前端传入的所查找的bookId页压缩进Redis中
func GetRedisBookId(c *gin.Context, bookIdStr string, size int, direction string) []Book {
	key := "bookId" + bookIdStr //为查询页
	id, err := strconv.Atoi(bookIdStr)
	fmt.Println(key)
	//根据传入的pageStr构造键名key，用于在Redis中查找数据,并将数据保存在data中.
	data, err := global.RedisConn.Get(c, key).Bytes()
	//如果Redis中不存在该数据，即err为redis.Nil，则说明需要从数据库中获取数据，并进行压缩存储到Redis中。
	if err == redis.Nil {
		book := make([]Book, 0)
		if direction == "back" {
			sql := "select * from book where id > ? order by id ASC limit ?"
			err := global.DB.Raw(sql, id, size).Scan(&book).Error
			if err != nil {
				fmt.Printf("查找bookId页之后的图书失败:%+v\n", err.Error())
			}
		} else {
			sql := "select * from book where id < ? order by id ASC limit ?"
			err := global.DB.Raw(sql, id, size).Scan(&book).Error
			if err != nil {
				fmt.Printf("查找bookId页之前的图书失败:%+v\n", err.Error())
			}
		}

		//进行压缩
		Y := YS2(book)
		data = Y
		//使用RedisConn.Set方法将key和Y存储到Redis中，并设置过期时间为5分钟。
		num := rand.Intn(5)
		err = global.RedisConn.Set(c, key, Y, time.Duration(num)*time.Second).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			fmt.Println("压缩错误 DB2：Y")
			return nil
		}
	}
	//调用SelectBooks函数对data进行解压和反序列化，将解压后的数据返回。
	//序列化（Serialization）是指将数据结构或对象状态转换为可存储或传输的格式的过程，例如将对象转换为字节序列或字符序列。
	//反序列化（Deserialization）则是将序列化的数据重新转换为原始数据结构或对象的过程。
	return SelectBooks(data)
}

func PreHeating(bookId int64, size int) []Book {
	var books []Book
	sql := "select * from book where id >= ? limit ?"
	err := global.DB.Raw(sql, bookId, size).Scan(&books).Error
	if err != nil {
		fmt.Printf("获取数据库中热门书籍失败:err:%+v\n", err.Error())
	}
	return books
}

func SavePreHeatingBooks(id int64, book []Book) {
	c := context.Background()
	value := YS(book)
	idStr := strconv.FormatInt(id, 10)
	key := "book" + idStr
	err := global.RedisConn.Set(c, key, value, 5*time.Second).Err()
	if err != nil {
		fmt.Printf("将预热数据存储在redis中失败！err:%+v\n", err.Error())
	}
}

// GetBook 已弃用
func GetBook(bookId int64) []*Book {
	book := make([]*Book, 0)
	//var book=&Book{}
	//sql"="SELECT *FROM `book` WHERE id =" +string(bookId)
	//sql := fmt.Sprintf("SELECT *FROM `book` WHERE id = %d", bookId)
	sql := "select * from book where  id = ?"
	err := DB.Raw(sql, bookId).Scan(&book).Error
	//err := DB.Table("book").Where("id = ?", bookId).Scan(&book).Error
	if err != nil {
		fmt.Printf("查找id为%d的图书失败！\n", bookId)
	}
	return book
}

func AddBooks(categoryId int64, bn, name, description string, count int) bool {
	//var book *Book
	tx := DB.Begin()
	sql := "INSERT INTO book (category_id, bn, name,description,count) VALUES (?,?,?,?,?)"
	err := tx.Exec(sql, categoryId, bn, name, description, count).Error
	if err != nil {
		fmt.Printf("添加图书插入数据失败！ err:%+v\n", err.Error())
		tx.Rollback()
	}
	fmt.Println("添加图书信息成功！")
	tx.Commit()
	return true
}

func UpdateBooks(bookId int64, bn, name, description string, count int, categoryId int64) bool {
	tx := DB.Begin()
	sql := "UPDATE book SET bn=?, name=?, description=?, count=?, category_id=? WHERE id=?"
	err := tx.Exec(sql, bn, name, description, count, categoryId, bookId).Error
	if err != nil {
		fmt.Printf("更新图书信息失败！err：%+v\n", err.Error())
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func DeleteBook(bookId int64) bool {
	tx := DB.Begin()
	sql := "DELETE FROM book WHERE id = ?"
	err := tx.Exec(sql, bookId).Error
	if err != nil {
		fmt.Printf("删除用户信息失败！err:%+v\n", err.Error())
		tx.Rollback()
		return false
	}

	tx.Commit()
	return true
}

func AddCategory(name string) bool {
	tx := DB.Begin()
	sql := "INSERT INTO category (id, name) VALUES (?,?)"
	err := tx.Exec(sql, name).Error
	if err != nil {
		fmt.Printf("插入新图书种类失败！err:%+v\n", err.Error())
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func UpdateCategory(id int64, name string) bool {
	tx := DB.Begin()
	sql := "UPDATE category SET name = ? WHERE id = ?"
	err := tx.Exec(sql, name, id).Error
	if err != nil {
		fmt.Printf("更新图书种类失败！err:%+v\n", err.Error())
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func DeleteCategory(id int64) bool {
	tx := DB.Begin()
	sql := "DELETE FROM category WHERE id = ?"
	err := tx.Exec(sql, id).Error
	if err != nil {
		fmt.Printf("删除图书种类信息失败！err:%+v\n", err.Error())
		tx.Rollback()
		return false
	}

	tx.Commit()
	return true
}
