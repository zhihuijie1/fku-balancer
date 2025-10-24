package midWare

import (
	"fku-balancer/config"
	"fmt"
	"net/http"
)

func PathAnalysMidWare(config *config.Config) func(http.Handler) http.Handler {
	pathMap := make(map[string]bool)
	for _, l := range config.Location {
		pathMap[l.Pattern] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !pathMap[r.URL.Path] {
				fmt.Println("PathAnalysMidWare - parten匹配失败")
				http.NotFound(w, r)
				return
			}
			fmt.Println("PathAnalysMidWare - parten匹配成功进入一下环节")
			next.ServeHTTP(w, r)
			fmt.Println("PathAnalysMidWare - return")
		})
	}

}
