package main

import (
	"github.com/Tansuozhe1num/codedream/internal/preload"
	"github.com/Tansuozhe1num/codedream/internal/router"
	"log"
)

func main() {
	db, err := preload.InitDB()
	if err != nil {
		log.Fatalf("init DB failed")
		// 错误了还是咔嚓掉
		return
	}

	// // 配置cornjob： 每天爬取部分题目
	//if err = preload.StartCronJob(db); err != nil {
	//	log.Fatalf("init DB failed")
	//	return
	//}

	rt := router.NewRouter(db)
	err = rt.Run(":8080")
	if err != nil {
		log.Fatalf("init DB failed")
	}
}
