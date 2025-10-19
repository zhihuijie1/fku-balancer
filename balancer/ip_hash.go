// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// IP哈希负载均衡算法实现
// 根据客户端IP地址选择固定的服务器，实现会话保持
package balancer

import (
	// hash/crc32 提供CRC32校验和算法
	// CRC（Cyclic Redundancy Check）循环冗余校验
	// 用于将IP地址字符串转换为数字，实现哈希功能
	"hash/crc32"
)

// init函数在包被导入时自动执行
// 注册IP哈希算法到工厂映射表
func init() {
	// 将IP哈希算法注册到全局工厂映射表
	// IPHashBalancer 是算法名称常量
	// NewIPHash 是创建IPHash负载均衡器的工厂函数
	factories[IPHashBalancer] = NewIPHash
}

// IPHash IP哈希负载均衡器结构体
// 基于客户端IP地址选择服务器，实现会话保持（Session Affinity）
//
// 算法原理：
// 1. 将客户端IP地址通过哈希函数转换为一个数值
// 2. 用这个数值对服务器数量取模，得到服务器索引
// 3. 同一个IP地址总是会被哈希到同一个服务器
//
// 优点：
// 1. 会话保持：同一客户端的请求总是发送到同一服务器
// 2. 无状态：不需要在负载均衡器中保存会话信息
// 3. 实现简单，性能好
// 4. 适合需要保持用户状态的应用
//
// 缺点：
// 1. 负载可能不均：如果某些IP的请求量特别大，会导致负载不均
// 2. 服务器变化影响大：增加或减少服务器会导致大量会话重新分配
// 3. 不考虑服务器的实际负载和处理能力
// 4. NAT环境下效果不好（多个用户可能共享同一个公网IP）
type IPHash struct {
	// 嵌入BaseBalancer，继承基础功能
	// 包括Add、Remove方法和线程安全的hosts管理
	// 这是Go的组合模式，通过嵌入实现代码复用
	BaseBalancer
}

// NewIPHash 创建一个新的IP哈希负载均衡器
// 这是工厂函数，实现了Factory函数签名
//
// 参数 hosts：初始的服务器地址列表
// 返回值：实现了Balancer接口的IPHash实例
//
// 设计说明：
// IPHash结构体很简单，只需要继承BaseBalancer
// 核心逻辑在Balance方法中实现
func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		// 初始化基础负载均衡器
		// 设置服务器列表
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

// Balance 根据IP哈希算法选择服务器
// 实现了Balancer接口的核心方法
//
// 参数 key：客户端标识，通常是IP地址
//
//	这个参数在IPHash算法中非常重要
//
// 返回值：
//   - string: 选中的服务器地址
//   - error: 可能的错误（如没有可用服务器）
//
// 算法流程：
// 1. 获取读锁，保护hosts切片
// 2. 检查是否有可用服务器
// 3. 计算IP地址的哈希值
// 4. 取模得到服务器索引
// 5. 返回对应的服务器
func (r *IPHash) Balance(key string) (string, error) {
	// 获取读锁
	// 保护hosts切片的并发读取
	r.RLock()
	// defer确保函数返回前释放锁
	defer r.RUnlock()

	// 检查是否有可用的服务器
	if len(r.hosts) == 0 {
		// 没有可用服务器，返回错误
		return "", NoHostError
	}

	// 核心算法：IP哈希
	// crc32.ChecksumIEEE 使用IEEE多项式计算CRC32校验和
	// 为什么使用CRC32？
	// 1. 计算速度快，性能好
	// 2. 分布均匀，减少哈希冲突
	// 3. 标准算法，结果可预测
	//
	// []byte(key) 将字符串转换为字节数组
	// Go中字符串是不可变的，但可以高效转换为[]byte
	//
	// % uint32(len(r.hosts)) 取模运算
	// 将哈希值映射到[0, len(hosts))范围内
	// 使用uint32确保类型匹配
	//
	// 数学原理：
	// 哈希函数：H(key) = CRC32(key)
	// 服务器索引：index = H(key) mod N
	// 其中N是服务器数量
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.hosts))

	// 返回选中的服务器
	// value是uint32类型，可以安全地用作数组索引
	return r.hosts[value], nil
}

// IP哈希算法的应用场景：
// 1. 购物车系统：用户的购物车数据保存在特定服务器上
// 2. 用户会话：保持用户登录状态
// 3. 游戏服务：玩家状态需要保持在同一服务器
// 4. WebSocket连接：长连接需要固定服务器
//
// 不适用场景：
// 1. 无状态服务：不需要会话保持
// 2. 动态扩缩容：服务器经常变化
// 3. 热点问题：某些IP请求量特别大
//
// 性能分析：
// - 时间复杂度：O(1)
// - 空间复杂度：O(n)，n是服务器数量
// - CRC32计算：O(m)，m是key的长度，通常很小
//
// 哈希算法的选择：
// 1. CRC32：快速，本项目使用
// 2. MD5：更均匀，但较慢
// 3. MurmurHash：性能和分布都很好
// 4. xxHash：极快，适合高性能场景
//
// 一致性问题：
// 当服务器数量变化时，大部分请求会被重新分配
// 解决方案：使用一致性哈希算法（ConsistentHash）
//
// 热点问题的解决方案：
// 1. 虚拟节点：一个物理服务器映射多个虚拟节点
// 2. 请求分片：将热点IP的请求分散到多个服务器
// 3. 动态权重：根据负载动态调整服务器权重
//
// 安全性考虑：
// 如果攻击者知道哈希算法，可能构造特定IP导致负载不均
// 解决方案：添加随机盐值，或使用加密哈希函数
