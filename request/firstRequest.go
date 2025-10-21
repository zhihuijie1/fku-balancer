package request

import (
	"io"
	"net/http"
)

func FirstRequest() {
	// 忽略错误，仅用于测试
	resp, err := http.Get("http://localhost:8088")
	if err != nil {
		println("请求失败:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("读取响应体失败:", err.Error())
		return
	}
	println(resp.Status)
	println("请求成功")
	println(string(body))
}
