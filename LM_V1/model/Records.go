package model

import (
	"LibraryManagementV1/LM_V1/tools"
	"fmt"
	"time"
)

// CreateRecord 用户借阅图书记录
func CreateRecord(userId int64, bookId int64) *Record {
	var ret Book

	tx := DB.Begin()

	sql := "SELECT * FROM book WHERE id = ?"
	err := tx.Raw(sql, bookId).Find(&ret).Error
	if err != nil {
		tx.Rollback()
	}

	if ret.Id <= 0 {
		tx.Rollback()
		return nil
	}
	fmt.Println(ret)

	record := &Record{
		UserId:    userId,
		BookId:    bookId,
		Status:    1,
		StartTime: time.Now(),
		OverTime:  time.Now().Add(tools.ContinueTime),
	}
	sql = "INSERT INTO record (user_id,book_id,status,start_time,over_time) VALUES (?,?,?,?,?)"
	err = tx.Exec(sql, record.UserId, record.BookId, record.Status, record.StartTime, record.OverTime).Error
	if err != nil {
		tx.Rollback()
	}
	sql = "update book set count=count-1 where id = ?"
	err = tx.Exec(sql, record.BookId).Error
	if err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return record
}

func ReturnBook(userId int64, bookId int64) {
	tx := DB.Begin()
	tm := time.Now()
	var book Book
	sql := "SELECT  *FROM `book` WHERE id = ?"
	err := tx.Raw(sql, bookId).Find(&book).Error
	if err != nil {
		fmt.Printf("err:%+v\n", err.Error())
		tx.Rollback()
		return
	}
	if book.Id <= 0 {
		tx.Rollback()
		return
	}

	sql = "UPDATE record SET status = 0, return_time = ? WHERE book_id = ? AND user_id = ?"
	err = tx.Exec(sql, tm, bookId, userId).Error

	/*err := tx.Table("record").Where("book_id = ? and user_id = ?", bookId, userId).UpdateColumns(map[string]interface{}{
		"status":      0,
		"return_time": tm,
	}).Error*/
	if book.Id <= 0 {
		fmt.Println("222222")
		tx.Rollback()
		return
	}
	sql = "UPDATE book SET count=count+1 WHERE id = ?"
	err = tx.Exec(sql, bookId).Error
	//err = tx.Table("book").Where("id = ?", bookId).Update("count", gorm.Expr("count + ?", 1)).Error
	tx.Commit()
	return
}

func UserGetRecords(userId int64) []*Record {
	var record []*Record
	sql := "select * from record where user_id = ? order by start_time"
	err := DB.Raw(sql, userId).Scan(&record).Error
	if err != nil {
		fmt.Printf("查询用户借阅记录失败！err:%+v\n", err.Error())
	}
	return record
}

func AdminGetRecords() *Record {
	var record *Record
	sql := "select * from record "
	err := DB.Raw(sql).Scan(&record).Error

	if err != nil {
		fmt.Printf("管理员查询借阅记录表失败！err:%+v\n", err.Error())
		return nil
	}
	return record
}

func GetUserRecordStatus(status int) *Record {
	var record *Record
	sql := "select * from record where status = ?"
	err := DB.Raw(sql, status).Scan(&record).Error
	if err != nil {
		fmt.Printf("管理员查询借阅归还或未归还状态失败！err:%+v\n", err.Error())
		return nil
	}
	return record
}
