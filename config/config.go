package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// 负责读取、解析和使用变量存储配置文件中自定义的配置
// SSLCertificateKey 当schema为https时,存储https的私钥文件路径
// SSLCertificate 当schema为https时,存储https的证书文件路径
type Config struct {
	Schema                string      `yaml:"schema"`
	Port                  int         `yaml:"port"`
	Tcp_health_check      bool        `yaml:"tcp_health_check"`
	Health_check_interval uint        `yaml:"health_check_interval"`
	Max_allowed           uint        `yaml:"max_allowed"`
	Location              []*Location `yaml:"location"`
	SSLCertificateKey     string      `yaml:"ssl_certificate_key"`
	SSLCertificate        string      `yaml:"ssl_certificate"`
}

type Location struct {
	Pattern      string   `yaml:"pattern"`
	Proxy_pass   []string `yaml:"proxy_pass"`
	Balance_mode string   `yaml:"balance_mode"`
}

func ReadConfig(filename string) (*Config, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *Config) Print() {
	fmt.Printf("schema: %s\n port: %s\n, tcp_health_check: %s\n, health_check_interval: %s\n, max_allowed: %s\n",
		config.Schema, config.Port, config.Tcp_health_check, config.Health_check_interval, config.Max_allowed)

	l := config.Location
	for _, v := range l {
		fmt.Printf("pattern: %s\n, proxy_pass: %s\n, palance_mode: %s\n", v.Pattern, v.Proxy_pass, v.Balance_mode)
	}
}

// 验证配置文件的合理性
func (c *Config) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("the schema \"%s\" not supported", c.Schema)
	}

	if c.Schema == "https" && (len(c.SSLCertificate) <= 0 || len(c.SSLCertificateKey) <= 0) {
		return errors.New("the https proxy requires ssl_certificate_key and ssl_certificate")
	}

	if c.Port <= 1 {
		fmt.Errorf("Port must be greater than 0")
	}

	if len(c.Location) <= 0 {
		return errors.New("the details of location cannot be null")
	}

	// 健康检查必须大于0.否则我写的这个程序就出bug
	if c.Health_check_interval < 1 {
		return errors.New("health_check_interval must be greater than 0")
	}

	return nil
}
