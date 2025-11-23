package db

import (
	"context"
	"demo-project/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init 初始化MySQL数据库连接
func Init(ctx context.Context, cfg *config.DatabaseConfig) error {
	if cfg == nil {
		return fmt.Errorf("数据库配置不能为空")
	}

	// 只支持MySQL
	if cfg.Driver != "mysql" {
		return fmt.Errorf("不支持的数据库类型：%s，只支持mysql", cfg.Driver)
	}

	// 建立GORM连接
	gormDB, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	// 测试连接
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	db = gormDB
	fmt.Printf("✅ MySQL连接成功 (连接池: 最大%d/空闲%d)\n", cfg.MaxOpenConns, cfg.MaxIdleConns)
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		panic("❌ 数据库未初始化，请先调用db.Init()")
	}
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取连接池失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	fmt.Println("✅ 数据库连接已关闭")
	return nil
}

// HealthCheck 健康检查
func HealthCheck(ctx context.Context) error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取连接池失败: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("数据库健康检查失败: %w", err)
	}

	return nil
}
