package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level   string `mapstructure:"level"`
	LogFile string `mapstructure:"log_file"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	Secret       string `mapstructure:"secret"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// AppConfig 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 重置Viper配置
	viper.Reset()

	// 设置默认值
	setDefaultConfig()

	// 从配置文件加载
	loadConfigFile()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置错误: %w", err)
	}

	// 设置全局配置
	AppConfig = &config

	// 打印配置信息（仅用于调试，非release模式）
	if config.Server.Mode != "release" {
		log.Printf("加载的配置: %+v", config)
	}
	return nil
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	// 服务器配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.secret", "chigua-secret-key")

	// 数据库配置
	viper.SetDefault("database.host", "postgres")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "chigua")
	viper.SetDefault("database.sslmode", "disable")

	// 日志配置
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.log_file", "logs/app.log")
}

// loadConfigFile 加载配置文件
func loadConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("读取配置文件错误: %v", err)
		} else {
			log.Println("未找到配置文件，使用默认值和环境变量")
		}
	} else {
		log.Println("配置文件读取成功")
	}
}
