package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig 应用配置
type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Dsn             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*AppConfig, error) {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析YAML
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	return &config, nil
}
