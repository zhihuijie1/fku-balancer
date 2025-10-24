package main

import (
	"fku-balancer/config"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

var (
	pathIsStart = make(map[string]bool)
	wg          sync.WaitGroup
)

// return: path:bool
func PathIsStart(config *config.Config) {

	for _, l := range config.Location {
		for _, path := range l.Proxy_pass {
			pathIsStart[path] = false
		}
	}
	startHost()
}

func startHost() {
	// 取端口
	portSlice := make([]string, 0)
	portSet := make(map[string]bool) // 用于去重

	for key := range pathIsStart {
		// 解析URL提取端口
		u, err := url.Parse(key)
		if err != nil {
			continue
		}

		// 获取端口，如果没有指定端口则使用默认端口
		port := u.Port()
		if port == "" && u.Scheme == "http" {
			port = "80"
		} else if port == "" && u.Scheme == "https" {
			port = "443"
		}

		// 避免重复启动相同端口
		if port != "" && !portSet[port] {
			portSlice = append(portSlice, port)
			portSet[port] = true
		}
	}

	for _, port := range portSlice {
		go func(port string) {
			fmt.Println("启动服务，端口：", port)
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				fmt.Printf("服务器在端口 %s 启动失败: %v\n", port, err)
			}
		}(port)
	}

}

func hostHandler(w http.ResponseWriter, r *http.Request) {
	// 处理主机请求的逻辑
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Host is running"))
}
