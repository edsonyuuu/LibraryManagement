package model

import (
	"LibraryManagementV1/LM_V4/global"
	"fmt"
	"time"
)

func CheckAdminMsg(userName, password string) *Librarian {

	var admin *Librarian
	sql := "select * from librarian where user_name = ? and password = ?"
	err := global.DB.Raw(sql, userName, password).Scan(&admin).Error
	if err != nil {
		fmt.Printf("管理员登录信息校验失败！")
	}
	return admin
}
func DeleteUser(userId int64) bool {
	tx := global.DB.Begin()
	sql := "delete from user where id = ? "
	err := tx.Exec(sql, userId).Error
	if err != nil {
		fmt.Printf("删除用户id=%d操作失败！\n", userId)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func ReturnTime(time1 time.Time) []int64 {
	user := make([]int64, 0)
	sql := "select distinct user_id from record where status = 1 and over_time <= ?"
	err := global.DB.Raw(sql, time1).Scan(&user).Error
	if err != nil {
		fmt.Printf("查询未归还用户时间err:%+v\n", err.Error())
	}
	return user
}

func AdvanceOneDay() {
	user := make([]*Record, 0)
	//user := make([]*int64, 0)
	sql := "select distinct user_id from record where status = 1 and  over_time-Now()<interval '1 day'"
	err := global.DB.Raw(sql).Scan(&user).Error
	if err != nil {
		fmt.Printf("查找到期前一天的用户id失败!err:%+v\n", err.Error())
	}
	return
}

func BanUsers(userId []int64) bool {
	user := make([]*Record, 0)
	sql1 := "select user_id from record where status = 1 "
	err := global.DB.Raw(sql1, userId).Scan(&user).Error
	if err != nil {
		fmt.Printf("sql查询失败!err:%+v\n", err.Error())
	}
	//
	tx := global.DB.Begin()
	sql2 := "update user set status = 1  where id = ?"
	for i := 0; i < len(userId); i++ {
		tx.Exec(sql2, userId[i])
	}
	err = tx.Exec(sql2, userId).Error
	if err != nil {
		fmt.Printf("封禁用户状态失败！err:%+v\n", err.Error())
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

//func WillReturnUser(time time.Time) {
//	sql := "select "
//}

/*func AdminGetUserBorrowStatus(userId int64, status int) *Record {
	var record *Record
	sql := "select * from record where user_id = ? and status = ?"
	err := DB.Raw(sql, userId, status).Scan(&record).Error
	if err != nil {
		fmt.Printf("管理员查询用户借阅状态失败！err：%+v\n", err.Error())
	}
	return record
}
*/
