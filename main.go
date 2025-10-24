package main

import (
	"fku-balancer/config"
	"fku-balancer/midWare"
	"fku-balancer/proxy"
	"fku-balancer/request"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// 1读配置文件
	config, err := config.ReadConfig("config/config.yaml")
	if err != nil {
		// log.Fatalf：打印错误信息后调用os.Exit(1)终止程序
		log.Fatalf("read config error: %s", err)
	}

	// 2配置文件合法性校验
	err = config.Validation()
	if err != nil {
		log.Fatalf("validation config error: %s", err)
	}

	PathIsStart(config)

	// 3 启动http路由器
	// gorilla/mux是一个功能强大的URL路由器和调度器
	router := mux.NewRouter()

	// 4为每个路由配置创建反向代理
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHttpProxy(l.Proxy_pass, l.Balance_mode)

		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}

		if config.Tcp_health_check {
			httpProxy.HealthCheck(config.Health_check_interval)
		}

		router.Handle(l.Pattern, httpProxy)
	}

	// 5添加中间件（如果配置了最大并发数）
	// 中间件是在请求到达处理器之前/之后执行的代码
	// 这里的中间件用于限制并发请求数
	// 传入的参数是: mwf ...MiddlewareFunc,这是一个可变参数,可以传入多个中间件
	// 当有多个中间件的时候,middlware会按照顺序执行,可以直接传入一个中间件切片
	// MiddlewareFunc 是一个函数类型 type MiddlewareFunc func(http.Handler) http.Handler
	midwares := []mux.MiddlewareFunc{
		midWare.PathAnalysMidWare(config),
		midWare.MaxRequestMidWare(config.Max_allowed),
	}
	for _, mid := range midwares {
		router.Use(mid)
	}

	// 第八步：创建HTTP服务器
	server := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	// 第九步：打印配置信息
	config.Print()

	// 第十步：启动服务器监听

	go func() {
		time.Sleep(2 * time.Second)

		request.FirstRequest()
	}()

	if config.Schema == "http" {
		server.ListenAndServe()
		if err != nil {
			// 通常是端口被占用或权限不足
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}
