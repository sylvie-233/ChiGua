package database

import (
	"chigua-backend/config"
	"chigua-backend/utils/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB 数据库连接实例
var DB *sqlx.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	// 构建连接字符串
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)

	logger.Infof("正在连接数据库: %s:%s/%s", config.AppConfig.Database.Host, config.AppConfig.Database.Port, config.AppConfig.Database.DBName)

	// 连接数据库
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("测试数据库连接失败: %w", err)
	}

	// 保存连接实例
	DB = db

	logger.Infof("数据库连接成功: %s:%s/%s", config.AppConfig.Database.Host, config.AppConfig.Database.Port, config.AppConfig.Database.DBName)

	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		logger.Info("正在关闭数据库连接")
		err := DB.Close()
		if err != nil {
			logger.Errorf("关闭数据库连接失败: %v", err)
			return err
		}
		logger.Info("数据库连接关闭成功")
	}
	return nil
}
