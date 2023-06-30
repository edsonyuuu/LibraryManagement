package model

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
)

func GetCategory() []*Category {
	category := make([]*Category, 0)
	sql1 := "SELECT * FROM `category`"
	err := DB.Raw(sql1).Scan(&category).Error

	if err != nil {
		fmt.Printf("查询图书种类err:%+v\n", err.Error())
	}
	return category
}

// GetBooks    []*Book
//books := make([]*Book, 0)
/*sql := "SELECT * FROM `book`"

err1 := DB.Raw(sql).Scan(&books).Error
if err1 != nil {
	fmt.Printf("查询图书信息err:%+v\n", err1.Error())
}
return books*/
func GetBooks() []Book {
	book, _ := SelectBooks()
	if book != nil {
		return book
	}
	return nil

}

func SelectBooks() ([]Book, error) {
	// 压缩过暂时注释掉
	//rows, err := DB.Raw("SELECT * FROM book").Rows()
	//YS(rows)
	ctx := context.Background() // 创建上下文
	// 查询以 "book" 为前缀的键
	// 查询 JSON 数据
	val, err := RedisConn.Get(ctx, "books").Bytes()
	if err != nil {
		fmt.Printf("查询redis 数据出现错误！err:%+v\n", err.Error())
		return nil, err
	}

	// 创建 gzip 解码器
	gzipReader, err := gzip.NewReader(bytes.NewReader(val))
	if err != nil {
		fmt.Printf("vgzip.NewReader 时出现错误！err:%+v\n", err.Error())
		return nil, err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			fmt.Printf("gzipReader.Close() 时出现错误！err:%+v\n", err.Error())
			return
		}
	}(gzipReader)
	//
	books := make([]Book, 0)
	// 分批读取解压缩后的 JSON 数据
	decoder := json.NewDecoder(gzipReader)
	batchSize := 10 // 每次读取的批次大小为 100 条记录
	for {
		var batch []Book
		err := decoder.Decode(&batch)
		if err == io.EOF { // 已经读取完数据
			break
		} else if err != nil {
			fmt.Printf("ecoder.Decode(&batch) 时出现错误！err:%+v\n", err.Error())
			return nil, err
		}
		books = append(books, batch...)
		if len(books) >= batchSize { // 达到批次大小，返回结果
			break
		}
	}
	return books, nil
}

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

func AddBooks(categoryId int64, bn, name, description string, count int) *Book {
	var book *Book
	tx := DB.Begin()
	sql := "INSERT INTO book (category_id, bn, name,description,count) VALUES (?,?,?,?,?)"
	err := tx.Exec(sql, categoryId, bn, name, description, count).Error
	if err != nil {
		fmt.Printf("添加图书插入数据失败！ err:%+v\n", err.Error())
		tx.Rollback()
	}
	tx.Commit()
	return book
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
