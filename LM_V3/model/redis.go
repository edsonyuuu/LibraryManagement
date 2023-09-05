package model

import (
	"LibraryManagementV1/LM_V3/global"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisConn 存储全部图书信息
//var RedisConn *redis.Client

// RedisConn2 存储查找bookId所在页
//var RedisConn2 *redis.Client

func InitRedis() {
	// 创建 Redis 客户端连接
	global.RedisConn = redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr(),
		Password: global.Config.Redis.Password, // Redis 未设置密码时为空
		DB:       1,
	})
	// 测试连接是否成功
	_, err := global.RedisConn.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis DB2: %v", err)
		return
	}
	fmt.Println("Connected to Redis DB2")
}

// YS 压缩图书数据，以json形式压缩，但Redis中存储的类型还是为string。
func YS(book []Book) []byte {
	ctx := context.Background() // 创建上下文，用于在Redis连接中使用

	//如果Redis操作的响应时间过长，可能会导致应用程序的性能下降，甚至阻塞整个应用程序。
	//为了避免这种情况的发生，我们可以使用上下文来设置Redis操作的超时时间，一旦超时时间到达，就可以自动取消Redis操作，从而避免应用程序的阻塞。
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// 创建 gzip 编码器
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf) ////将编码器与buf关联起来

	// 将 books 切片转换为 JSON 数据，并压缩保存在jsonBytes
	jsonBytes, err := json.Marshal(&book)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) :%+v\n", err.Error())
	}
	// 使用gzipWriter将jsonBytes进行压缩
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write :%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close :%+v\n", err.Error())

	}

	// 将压缩后的 JSON 数据存储到 Redis 中
	// 使用RedisConn.Set方法将buf.Bytes()作为值，"book"作为键，存储到Redis中，并设置过期时间为5分钟。
	err = global.RedisConn.Set(ctx, "book", buf.Bytes(), 5*time.Minute).Err()

	//如果存储过程中出现错误，会在控制台输出相关错误信息。
	if err != nil {
		fmt.Printf("Error during set:%+v\n", err.Error())
	}
	//
	fmt.Println("Data imported to Redis successfully.")

	//返回buf.Bytes()，即压缩后的JSON数据，以便在需要时进行使用
	return buf.Bytes()
}

func YS2(book []Book) []byte {
	ctx := context.Background() // 创建上下文，用于在Redis连接中使用

	//如果Redis操作的响应时间过长，可能会导致应用程序的性能下降，甚至阻塞整个应用程序。
	//为了避免这种情况的发生，我们可以使用上下文来设置Redis操作的超时时间，一旦超时时间到达，就可以自动取消Redis操作，从而避免应用程序的阻塞。
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// 创建 gzip 编码器
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf) ////将编码器与buf关联起来

	// 将 books 切片转换为 JSON 数据，并压缩保存在jsonBytes
	jsonBytes, err := json.Marshal(&book)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) 2 :%+v\n", err.Error())
	}
	// 使用gzipWriter将jsonBytes进行压缩
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write 2:%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close 2:%+v\n", err.Error())

	}

	// 将压缩后的 JSON 数据存储到 Redis 中
	// 使用RedisConn.Set方法将buf.Bytes()作为值，"bookId"作为键，存储到Redis中，并设置过期时间为5分钟。
	err = global.RedisConn.Set(ctx, "bookId", buf.Bytes(), 5*time.Minute).Err()

	//如果存储过程中出现错误，会在控制台输出相关错误信息。
	if err != nil {
		fmt.Printf("Error during set DB2:%+v\n", err.Error())
	}
	//
	fmt.Println("Data imported to Redis successfully. DB2")

	//返回buf.Bytes()，即压缩后的JSON数据，以便在需要时进行使用
	return buf.Bytes()
}
