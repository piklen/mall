package main

import (
	"mall/cache"
	"mall/conf"
	"mall/mq"
	"mall/routes"
)

func main() {
	conf.Init()
	cache.InitCache()
	mq.InitRabbitMQ()
	r := routes.NewRouter()
	_ = r.Run(":7999")
}
