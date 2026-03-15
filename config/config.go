package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 是整个应用的配置根节点
type Config struct {
	MySQL   MySQLConfig   `yaml:"mysql"`
	Redis   RedisConfig   `yaml:"redis"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Gateway GatewayConfig `yaml:"gateway"`
	Chat    ChatConfig    `yaml:"chat"`
	RPC     RPCConfig     `yaml:"rpc"`
	JWT     JWTConfig     `yaml:"jwt"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

// DSN 生成 GORM 连接用的 Data Source Name
func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.DB)
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type GatewayConfig struct {
	HTTPPort int `yaml:"http_port"`
}

type ChatConfig struct {
	RPCPort int `yaml:"rpc_port"`
}

type RPCConfig struct {
	Framework string `yaml:"framework"` // minirpc | grpc | kitex
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // 小时
}

// Load 从指定路径加载 YAML 配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
