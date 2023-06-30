package model

/*import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
)*/
import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"time"
)

var store, _ = redis.NewStore(10, "tcp", "114.115.200.190:6379", "qwert", []byte("secret"))

// var store = sessions.NewCookieStore([]byte("secret"))
var sessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {

	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		fmt.Printf("GetsessionErr:%+v\n", err.Error())
	}
	fmt.Printf("session:%+v\n", session.Values)

	return session.Values
}

// SaveSession  设置保存session
func SaveSession(c *gin.Context, name string, id int64) error {
	//获取一个session对象
	store.Options(sessions.Options{MaxAge: int(24 * time.Hour)})
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		return err
	}

	//在session中存储值
	session.Values["name"] = name
	session.Values["id"] = id
	//保存更改
	return session.Save(c.Request, c.Writer)

}

// DeleteSession 删除session
func DeleteSession(c *gin.Context) error {
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		return err
	}
	fmt.Printf("session:%+v\n", session.Values)
	//session.Values["name"] = ""
	//session.Values["id"] = 0
	session.Options.MaxAge = -1 //这行实现将session删除
	return session.Save(c.Request, c.Writer)
}
