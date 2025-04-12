package database

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/funlake/wschat/pkg/orm"
	"github.com/spf13/viper"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// GetDB returns the database connection
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		initDB()
	})
	return db
}

// initDB initializes the database connection
func initDB() {
	var err error
	
	// 获取配置文件中的数据库参数
	dbType := viper.GetString("database.type")
	dbPath := viper.GetString("database.path")
	
	// 如果配置中没有设置，则使用默认值
	if dbPath == "" {
		dbPath = filepath.Join("data", "wschat.db")
		log.Printf("Database path not configured, using default: %s", dbPath)
	}
	
	log.Printf("Initializing %s database at: %s", dbType, dbPath)
	
	// 根据数据库类型连接数据库
	switch dbType {
	case "sqlite", "":
		// 确保目录存在
		dir := filepath.Dir(dbPath)
		if dir != "" && dir != "." {
			// 创建目录（如果不存在）
			if err := createDirIfNotExist(dir); err != nil {
				log.Fatalf("Failed to create database directory: %v", err)
			}
		}
		
		// 连接SQLite数据库
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	
	// 设置连接池参数
	maxIdleConns := viper.GetInt("database.max_idle_conns")
	maxOpenConns := viper.GetInt("database.max_open_conns")
	connMaxLifetime := viper.GetInt("database.conn_max_lifetime")
	
	if maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	if maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if connMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	}

	// 自动迁移数据库表结构
	log.Println("Auto-migrating database schema...")
	err = db.AutoMigrate(
		&orm.Customer{},
		&orm.Message{},
		&orm.Feedback{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	log.Println("Database initialized successfully")
}

// createDirIfNotExist 创建目录（如果不存在）
func createDirIfNotExist(dir string) error {
	// 使用os包的MkdirAll函数创建目录
	// 这个函数会创建所有不存在的父目录
	// 如果目录已经存在，则不会返回错误
	return os.MkdirAll(dir, 0755)
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting SQL DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}
} 