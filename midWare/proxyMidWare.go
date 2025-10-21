package midWare

import (
	"fmt"
	"net/http"
)

// 最大请求数中间件
func MaxRequestMidWare(maxReq uint) func(http.Handler) http.Handler {
	channel := make(chan struct{}, maxReq)
	add := func() {
		channel <- struct{}{}
	}

	remove := func() {
		<-channel
	}
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 当队列已满时，请求会被阻塞，直到队列有空闲位置
			add()
			defer remove()
			fmt.Println("proxyMidWare - 当前请求数：", len(channel), "最大请求数：", maxReq)
			next.ServeHTTP(w, r)
			fmt.Println("proxyMidWare - return")
		})
	}
}
