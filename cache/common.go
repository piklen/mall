package cache

import (
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"strconv"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// InitCache 在中间件中初始化redis链接  防止循环导包，所以放在这里
func InitCache() {
	Redis()
}

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint("2", 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       int(db),
	})
	_, err := client.Ping().Result() //心跳检测
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}
