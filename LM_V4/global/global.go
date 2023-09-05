package global

import (
	"LibraryManagementV1/LM_V4/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config    *config.Config
	DB        *gorm.DB
	Log       *logrus.Logger
	RedisConn *redis.Client
)
