package midWare

import (
	"fku-balancer/config"
	"fmt"
	"net/http"
)

// 路径分析中间件
func PathAnalysMidWare(config *config.Config) func(http.Handler) http.Handler {
	l := config.Location
	// 存储所有的路径parten
	partenTypeMap := make(map[string]bool)

	for _, v := range l {
		partenTypeMap[v.Pattern] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if partenTypeMap[r.URL.Path] {
				fmt.Println("PathAnalysMidWare - parten匹配成功进入一下环节")
				next.ServeHTTP(w, r)
				fmt.Println("PathAnalysMidWare - return")
			} else {
				fmt.Println("PathAnalysMidWare - parten匹配失败")
				http.NotFound(w, r)
			}
		})
	}
}
