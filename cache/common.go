package cache

import (
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

// RedisClient Redis缓存客户端单例
//var RedisClient *redis.Client

//type RedisClient struct {
//	*redis.Client
//	mu sync.Mutex
//}
//
//var redisClient *RedisClient
//
//// InitCache 在中间件中初始化redis链接  防止循环导包，所以放在这里
//func InitCache() {
//	Redis()
//}
//
//// Redis 在中间件中初始化redis链接
//func Redis() {
//	db, _ := strconv.ParseUint("2", 10, 64)
//	client := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "",
//		DB:       int(db),
//	})
//	_, err := client.Ping().Result() //心跳检测
//	if err != nil {
//		logging.Info(err)
//		panic(err)
//	}
//	redisClient = client
//}

// RedisClient 包含一个 *redis.Client 成员和一个互斥锁
type redisClient struct {
	Client *redis.Client
	Mu     sync.Mutex
}

var RedisClient *redisClient

// InitCache 在中间件中初始化 Redis 连接，防止循环导包，所以放在这里
func InitCache() {
	Redis()
}

// Redis 在中间件中初始化 Redis 连接
func Redis() {
	db, _ := strconv.ParseUint("2", 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       int(db),
	})
	_, err := client.Ping().Result() // 心跳检测
	if err != nil {
		logging.Info(err)
		panic(err)
	}

	// 初始化 RedisClient
	RedisClient = &redisClient{
		Client: client,
	}
}
