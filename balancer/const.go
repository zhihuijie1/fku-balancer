// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 负载均衡器算法名称常量定义
// 这个文件定义了所有支持的负载均衡算法名称
package balancer

// 定义负载均衡算法名称常量
// 使用const定义常量组，这些常量在编译时就确定值
// 常量的优点：
// 1. 类型安全：编译时检查
// 2. 性能更好：不占用内存，编译时替换
// 3. 不可修改：避免意外修改
//
// 这些常量用作工厂模式的key，用于创建对应的负载均衡器实例
// 在配置文件中，用户可以使用这些字符串来指定负载均衡算法
const (
	// IPHashBalancer IP哈希算法
	// 根据客户端IP地址计算哈希值，选择固定的服务器
	// 优点：同一客户端总是访问同一服务器，支持会话保持
	// 缺点：服务器增减时，会导致大量连接重新分配
	// 使用场景：需要会话保持的应用，如购物车、用户登录状态
	IPHashBalancer = "ip-hash"

	// ConsistentHashBalancer 一致性哈希算法
	// 将服务器和请求都映射到一个虚拟的环上，请求分配给环上顺时针方向最近的服务器
	// 优点：服务器增减时，只影响部分请求，最小化重新分配
	// 缺点：实现复杂，可能出现数据倾斜
	// 使用场景：分布式缓存、分布式存储
	ConsistentHashBalancer = "consistent-hash"

	// P2CBalancer Power of Two Choices（两次随机选择算法）
	// 随机选择两个服务器，然后选择负载较低的那个
	// 优点：比完全随机更均衡，比最少连接更高效
	// 缺点：需要维护连接数状态
	// 使用场景：高并发场景，需要在性能和负载均衡之间取得平衡
	// 算法来源：该算法被证明在理论和实践中都表现优秀
	P2CBalancer = "p2c"

	// RandomBalancer 随机算法
	// 随机选择一个服务器处理请求
	// 优点：实现简单，无状态，性能好
	// 缺点：可能出现负载不均，特别是请求量小的时候
	// 使用场景：服务器性能相近，请求处理时间相似
	RandomBalancer = "random"

	// R2Balancer Round-Robin（轮询算法）
	// 按照服务器列表顺序，依次分配请求
	// 优点：实现简单，保证每个服务器获得相同数量的请求
	// 缺点：不考虑服务器的实际负载和处理能力差异
	// 使用场景：服务器配置相同，请求处理时间相近
	// 命名说明：R2是Round-Robin的缩写
	R2Balancer = "round-robin"

	// LeastLoadBalancer 最少负载算法（也称最少连接算法）
	// 选择当前负载最小（活跃连接数最少）的服务器
	// 优点：考虑服务器实时负载，分配更合理
	// 缺点：需要维护连接数状态，实现较复杂
	// 使用场景：请求处理时间差异大，服务器性能不同
	// 实现说明：使用斐波那契堆实现，查找最小值O(1)，插入删除O(log n)
	LeastLoadBalancer = "least-load"

	// BoundedBalancer 有界一致性哈希算法
	// 一致性哈希的改进版，限制每个服务器的最大负载
	// 优点：结合一致性哈希和负载均衡的优点
	// 缺点：实现复杂，需要维护更多状态
	// 使用场景：需要一致性哈希特性，同时要避免热点问题
	// 算法说明：当某个服务器负载过高时，请求会分配给其他服务器
	BoundedBalancer = "bounded"
)

// 算法选择建议：
// 1. 简单场景：round-robin 或 random
// 2. 需要会话保持：ip-hash
// 3. 分布式缓存：consistent-hash
// 4. 高性能要求：p2c
// 5. 负载差异大：least-load
// 6. 综合场景：bounded
