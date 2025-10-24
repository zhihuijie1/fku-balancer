package proxy

import (
	"fku-balancer/balancer"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Http反向代理
// 反向代理是负载均衡器的核心组件，负责转发请求到后端服务器
var (
	ReverseProxy  = "Balancer-Reverse-Proxy"
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

// HTTPProxy 是HTTP反向代理的核心结构体
type HttpProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer
	alive   map[string]bool
	sync.RWMutex
}

// 把多个后端服务器地址转换成一个统一的HTTP代理，支持负载均衡和健康检查
func NewHttpProxy(targetHosts []string, algorithm string) (*HttpProxy, error) {
	hosts := make([]string, 0)
	hostsMap := make(map[string]*httputil.ReverseProxy)
	aliveMap := make(map[string]bool)

	for _, targetHost := range targetHosts {
		url, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}
		hostProxy := httputil.NewSingleHostReverseProxy(url)

		// 对发向后端服务器的请求进行改写
		originalDirector := hostProxy.Director
		hostProxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}

		host := GetHost(url)
		hosts = append(hosts, host)
		hostsMap[host] = hostProxy
		aliveMap[host] = true
	}

	lb, err := balancer.Build(algorithm, hosts)

	if err != nil {
		return nil, err
	}

	return &HttpProxy{
		hostMap: hostsMap,
		lb:      lb,
		alive:   aliveMap,
	}, nil
}

// ServeHTTP 实现http.Handler接口，处理HTTP请求
// 这是反向代理的核心方法，负责接收请求并转发
func (h *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	host, err := h.lb.Balance(GetIP(r))

	fmt.Println("拿到的host: ", host)

	if err != nil {
		// 选择失败（如没有可用服务器）
		// 返回502错误
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}

	h.lb.Inc(host)

	defer h.lb.Done(host)

	h.hostMap[host].ServeHTTP(w, r)
}
