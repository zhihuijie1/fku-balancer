package main

import (
	"github.com/gorilla/mux"
	"log"
)

func main() {
	// 1读配置文件
	config, err := ReadConfig("config.yml")
	if err != nil {
		// log.Fatalf：打印错误信息后调用os.Exit(1)终止程序
		log.Fatalf("read config error: %s", err)
	}

	// 2配置文件合法性校验
	err = config.Validation()
	if err != nil {
		log.Fatalf("validation config error: %s", err)
	}

	// 启动http路由器
	// gorilla/mux是一个功能强大的URL路由器和调度器
	router := mux.NewRouter()

}
