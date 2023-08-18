package model

import (
	"fmt"
	"time"
)

func Listen() {
	//
	sql := "select distinct user_id from record where over_time >=DATE_SUB(NOW(),INTERVAL 1 DAY) and status = 1"
	var userIds []int64
	DB.Raw(sql).Scan(&userIds)
	fmt.Println("所需还书用户id", userIds)

	if len(userIds) == 0 {
		return
	}
	//
	for _, userId := range userIds {
		tx := DB.Begin()
		//再次检测用户是否还书
		sql = "select * from record where status = 1 and user_id >= ? for update"
		var users []User
		tx.Raw(sql, userId).Scan(&users)

		//检测该用户最近30天是否已生成过未读消息
		sql = "select * from send_msg where status=1 and user_id=? and create_time>= DATE_SUB(CURDATE(), INTERVAL 30 DAY)"
		var message []SendMsg
		tx.Raw(sql, userId).Scan(&message)
		if len(users) != 0 && len(message) == 0 {
			sql = "insert into send_msg (user_id,msg,status,create_time) values(?,?,?,?)"
			err := tx.Exec(sql, userId, "您有书需要归还", 1, time.Now()).Error
			if err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				return
			}
		}
		tx.Commit()
	}
}
