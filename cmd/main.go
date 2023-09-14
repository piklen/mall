package main

import (
	"mall/cache"
	"mall/conf"
	"mall/routes"
)

func main() {
	conf.Init()
	cache.InitCache()
	r := routes.NewRouter()
	_ = r.Run(":7999")
}
