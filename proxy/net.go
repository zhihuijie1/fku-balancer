package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ConnectionTimeout = 3 * time.Second

// 根据请求,拿到客户端真实IP
func GetIP(r *http.Request) string {
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if len(r.Header.Get(XForwardedFor)) != 0 {
		xff := r.Header.Get(XForwardedFor)
		s := strings.Index(xff, ", ")
		if s == -1 {
			s = len(r.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(r.Header.Get(XRealIP)) != 0 {
		clientIP = r.Header.Get(XRealIP)
	}

	return clientIP
}

// 获取主机地址（不包含协议和路径）
// 例如：从"http://192.168.1.1:8080/path"提取"192.168.1.1:8080"
func GetHost(url *url.URL) string {
	if _, _, err := net.SplitHostPort(url.Host); err == nil {
		return url.Host
	}
	if url.Scheme == "http" {
		return fmt.Sprintf("%s:%s", url.Host, "80")
	} else if url.Scheme == "https" {
		return fmt.Sprintf("%s:%s", url.Host, "443")
	}
	return url.Host
}

// IsBackendAlive 检查后端服务是否存活
// 通过尝试建立 TCP 连接来判断目标主机是否可访问
func IsBackendAlive(host string) bool {
	// 解析主机地址为 TCP 地址
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}

	resolveAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)

	conn, err := net.DialTimeout("tcp", resolveAddr, ConnectionTimeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}
