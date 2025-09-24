// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// balancer包实现了各种负载均衡算法
// 负载均衡是分布式系统的核心组件，用于将请求分发到多个后端服务器
// 本包提供了统一的接口和多种算法实现
package balancer

import (
	"errors" // 标准错误包，用于创建错误对象
)

// 定义包级别的错误变量
// 使用var定义的包级变量会在程序启动时初始化
// 这些错误是预定义的，可以在整个包中复用
var (
	// NoHostError 表示没有可用的后端服务器
	// 当所有服务器都宕机或者服务器列表为空时返回此错误
	NoHostError = errors.New("no host")

	// AlgorithmNotSupportedError 表示请求的负载均衡算法不支持
	// 当Build函数收到未注册的算法名称时返回此错误
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

// Balancer 接口定义了负载均衡器的核心行为
// 这是一个策略模式的应用：定义一族算法，封装起来，让它们可以互相替换
// 所有的负载均衡算法都必须实现这个接口
//
// 接口的设计原则：
// 1. 接口应该小而专注（Interface Segregation Principle）
// 2. 依赖于抽象而不是具体实现（Dependency Inversion Principle）
type Balancer interface {
	// Add 添加一个新的后端服务器到负载均衡池
	// 参数：服务器地址字符串（如 "http://192.168.1.100:8080"）
	// 使用场景：服务器恢复正常、动态扩容
	Add(string)

	// Remove 从负载均衡池中移除一个后端服务器
	// 参数：服务器地址字符串
	// 使用场景：服务器宕机、动态缩容
	Remove(string)

	// Balance 根据负载均衡算法选择一个后端服务器
	// 参数：客户端标识（通常是IP地址），某些算法会用到（如ip-hash）
	// 返回值：选中的服务器地址和可能的错误
	// 这是负载均衡器的核心方法
	Balance(string) (string, error)

	// Inc 增加指定服务器的活跃连接数
	// 参数：服务器地址
	// 使用场景：开始处理请求时调用，用于least-connection等算法
	Inc(string)

	// Done 减少指定服务器的活跃连接数
	// 参数：服务器地址
	// 使用场景：请求处理完成时调用
	// Inc和Done配对使用，追踪服务器的实时负载
	Done(string)
}

// Factory 是创建Balancer的工厂函数类型
// 这是工厂模式（Factory Pattern）的实现
// 工厂模式的优点：
// 1. 解耦对象的创建和使用
// 2. 便于扩展新的算法
// 3. 统一的创建接口
//
// 参数：[]string 是初始的服务器地址列表
// 返回值：Balancer 接口的具体实现
type Factory func([]string) Balancer

// factories 是算法名到工厂函数的映射表
// 使用map存储所有注册的负载均衡算法
// key: 算法名称（如 "round-robin"、"random"）
// value: 对应的工厂函数
//
// 使用make初始化map，避免nil map panic
// 这个map会在各个算法的init函数中填充
var factories = make(map[string]Factory)

// Build 是负载均衡器的构建函数（Builder Pattern）
// 根据算法名称和服务器列表创建相应的负载均衡器实例
//
// 参数：
//   - algorithm: 算法名称，如 "round-robin"、"random"、"ip-hash" 等
//   - hosts: 初始的后端服务器地址列表
//
// 返回值：
//   - Balancer: 负载均衡器接口实例
//   - error: 如果算法不支持则返回错误
//
// 工作流程：
// 1. 从factories map中查找对应的工厂函数
// 2. 如果找到，调用工厂函数创建实例
// 3. 如果没找到，返回算法不支持错误
//
// 这种设计允许在运行时动态选择算法，实现了开闭原则：
// 对扩展开放（可以添加新算法），对修改关闭（不需要修改Build函数）
func Build(algorithm string, hosts []string) (Balancer, error) {
	// 从注册表中查找工厂函数
	// map的两值返回：value和是否存在
	factory, ok := factories[algorithm]
	if !ok {
		// 算法未注册，返回错误
		// 返回nil接口和错误是Go的惯用模式
		return nil, AlgorithmNotSupportedError
	}
	// 调用工厂函数创建负载均衡器实例
	// 工厂函数负责具体的初始化逻辑
	return factory(hosts), nil
}
