package model

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var RedisConn *redis.Client

func REdis() {

	// 创建 Redis 客户端连接
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Redis 未设置密码时为空
		DB:       1,  // 使用默认数据库
	})
	// 测试连接是否成功
	_, err := RedisConn.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v", err)
		return
	}
	fmt.Println("Connected to Redis")

}

func YS(sql *sql.Rows) {
	ctx := context.Background() // 创建上下文

	// 查询"book"表中的所有数据记录

	//rows, err := DB.Query("SELECT * FROM book")
	//if err != nil {
	//	fmt.Println("Failed to query book table:", err)
	//	return
	//}
	/*	defer func(sql *sql.Rows) {
			err := sql.Close()
			if err != nil {
				fmt.Printf("mysql close failed! err:%+v\n", err.Error())
			}
		}(sql)
	*/
	books := make([]Book, 0)
	// 遍历每条记录，并将其导入Redis
	for sql.Next() {
		book := Book{}
		// 读取记录中的字段值
		err := sql.Scan(&book.Id, &book.Name, &book.BN, &book.Description, &book.Count, &book.CategoryId, &book.ImgUrl)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		books = append(books, book)

	}

	// 检查遍历过程中是否有错误发生
	if err := sql.Err(); err != nil {
		fmt.Println("Error during iteration:", err)
		return
	}

	// 创建 gzip 编码器
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	// 将 books 切片转换为 JSON 数据，并压缩
	jsonBytes, err := json.Marshal(&books)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) :%+v\n", err.Error())
	}
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write :%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close :%+v\n", err.Error())

	}

	// 将压缩后的 JSON 数据存储到 Redis 中
	err = RedisConn.Set(ctx, "books", buf.Bytes(), 0).Err()
	if err != nil {
		fmt.Printf("Error during set:%+v\n", err.Error())
	}
	//

	fmt.Println("Data imported to Redis successfully.")
}

/*func init() {

// 连接到MySQL数据库
db, err := sql.Open("mysql", "yulongxin:rEjDCRjjdE7Fi5Sf@tcp(114.115.200.190:3306)/library")
if err != nil {
	fmt.Println("Failed to connect to MySQL:", err)
	return
}
defer func(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Printf("defer close mysql failed err:%+v\n", err.Error())
	}
}(db)
// 创建 Redis 客户端连接
RedisConn = redis.NewClient(&redis.Options{
	Addr:     "114.115.200.190:6379",
	Password: "qwert", // Redis 未设置密码时为空
	DB:       1,       // 使用默认数据库
})
// 测试连接是否成功
_, err = RedisConn.Ping(context.Background()).Result()
if err != nil {
	fmt.Printf("Failed to connect to Redis: %v", err)
	return
}
fmt.Println("Connected to Redis")

ctx := context.Background() // 创建上下文

// 查询"book"表中的所有数据记录
rows, err := db.Query("SELECT * FROM book")
if err != nil {
	fmt.Println("Failed to query book table:", err)
	return
}
defer func(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		fmt.Printf("mysql close failed! err:%+v\n", err.Error())
	}
}(rows)

books := make([]Book, 0)
// 遍历每条记录，并将其导入Redis
for rows.Next() {
	book := Book{}
	// 读取记录中的字段值
	err := rows.Scan(&book.Id, &book.Name, &book.BN, &book.Description, &book.Count, &book.CategoryId, &book.ImgUrl)
	if err != nil {
		fmt.Println("Failed to scan row:", err)
		continue
	}
	books = append(books, book)

	// 将记录导入Redis（示例为使用哈希表）
	/*	err = RedisConn.HSet(ctx, "book", fmt.Sprintf("%d", book.Id), fmt.Sprintf("%s-%s-%s-%d-%d-%s", book.Name, book.BN, book.Description, book.Count, book.CategoryId, book.ImgUrl)).Err()
		if err != nil {
			fmt.Println("Failed to set data in Redis:", err)
			continue
		}*/
//}

// 检查遍历过程中是否有错误发生
//if err := rows.Err(); err != nil {
//	fmt.Println("Error during iteration:", err)
//	return
//}
//data, err := json.Marshal(books)
//if err != nil {
//	panic(err)
//}
//
//err = RedisConn.Set(ctx, "books", data, 0).Err()
//if err != nil {
//	fmt.Printf("Error during set:%+v\n", err.Error())
//}
//
//fmt.Println("Data imported to Redis successfully.")
//}
