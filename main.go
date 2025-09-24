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

	// 3 启动http路由器
	// gorilla/mux是一个功能强大的URL路由器和调度器
	router := mux.NewRouter()

	// 4为每个路由配置创建反向代理

	// 5添加中间件（如果配置了最大并发数）
	// 中间件是在请求到达处理器之前/之后执行的代码
	// 这里的中间件用于限制并发请求数

	// 第八步：创建HTTP服务器

	// 第九步：打印配置信息

	// 第十步：启动服务器监听
}
