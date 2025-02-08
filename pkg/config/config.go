package config

import (
	"fmt"
	"os"
	"time"
	"gopkg.in/yaml.v3"
	"github.com/joho/godotenv"
)

// Config 系统配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig     `yaml:"jwt"`
	CORS     CORSConfig    `yaml:"cors"`
	Log      LogConfig     `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
	Issuer string       `yaml:"issuer"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins     []string `yaml:"allowed_origins"`
	AllowedMethods     []string `yaml:"allowed_methods"`
	AllowedHeaders     []string `yaml:"allowed_headers"`
	AllowCredentials   bool     `yaml:"allow_credentials"`
	MaxAge            int      `yaml:"max_age"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// GlobalConfig 全局配置实例
var GlobalConfig Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 1. 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果 .env 文件不存在，不返回错误，继续使用环境变量
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	// 2. 从环境变量加载配置
	loadFromEnv()

	// 3. 如果存在 YAML 配置文件，则从文件加载配置
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		// 尝试在项目根目录下查找配置文件
		configFile = "config.yaml"
		if _, err := os.Stat(configFile); err != nil {
			configFile = "../config.yaml"
		}
	}

	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %v", err)
		}

		if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
			return fmt.Errorf("failed to parse config file: %v", err)
		}

		// 解析时间配置
		if GlobalConfig.JWT.Expire != 0 {
			duration, err := time.ParseDuration(GlobalConfig.JWT.Expire.String())
			if err != nil {
				return fmt.Errorf("invalid JWT expire time: %v", err)
			}
			GlobalConfig.JWT.Expire = duration
		}
	}

	return nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv() {
	// 数据库配置
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		GlobalConfig.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		GlobalConfig.Database.Port = parseInt(dbPort, 3306)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		GlobalConfig.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		GlobalConfig.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		GlobalConfig.Database.DBName = dbName
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		GlobalConfig.JWT.Secret = secret
	}
	if expire := os.Getenv("JWT_EXPIRE"); expire != "" {
		if duration, err := time.ParseDuration(expire); err == nil {
			GlobalConfig.JWT.Expire = duration
		}
	}
	if issuer := os.Getenv("JWT_ISSUER"); issuer != "" {
		GlobalConfig.JWT.Issuer = issuer
	}

	// 服务器配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		GlobalConfig.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		GlobalConfig.Server.Port = parseInt(port, 8080)
	}
	if debug := os.Getenv("SERVER_DEBUG"); debug != "" {
		GlobalConfig.Server.Debug = debug == "true"
	}
}

// parseInt 解析整数，如果解析失败则返回默认值
func parseInt(s string, defaultValue int) int {
	var value int
	if _, err := fmt.Sscanf(s, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return &GlobalConfig
}
