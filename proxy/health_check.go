package proxy

import (
	"time"
)

func (h *HttpProxy) HealthCheck(interval uint) {
	// 对每一个服务器都要进行健康检查
	for host := range h.hostMap {
		go hostHealthCheck(host, h, interval)
	}
}

func hostHealthCheck(host string, h *HttpProxy, interval uint) {
	timeTicker := time.NewTicker(time.Duration(interval) * time.Second)
	for range timeTicker.C {
		if h.readAlive(host) && !IsBackendAlive(host) {
			h.setAlive(host, false)
			h.lb.Remove(host)
		} else if !h.readAlive(host) && IsBackendAlive(host) {
			h.setAlive(host, true)
			h.lb.Add(host)
		}
	}
}

// 读取后端服务器的存活状态
func (h *HttpProxy) readAlive(host string) bool {
	h.RLock()

	// 确保释放锁,即使发生panic,避免死锁
	defer h.RUnlock()

	return h.alive[host]
}

// 设置后端服务器的存活状态
func (h *HttpProxy) setAlive(host string, alive bool) {
	h.Lock()

	// 确保释放锁,即使发生panic,避免死锁
	defer h.Unlock()

	h.alive[host] = alive
}
