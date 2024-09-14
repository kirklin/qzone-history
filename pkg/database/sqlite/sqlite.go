package sqlite

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"qzone-history/pkg/database"
)

type SQLiteDB struct {
	db *gorm.DB
}

// NewSQLiteDB 创建一个新的 SQLite 数据库实例
// NewSQLiteDB creates a new SQLite database instance.
func NewSQLiteDB() database.Database {
	return &SQLiteDB{}
}

// Connect 连接到 SQLite 数据库
// Connect connects to the SQLite database.
// 如果 `DBName` 为空，则使用内存数据库。
// If `DBName` is empty, it will use an in-memory database.
func (s *SQLiteDB) Connect(config *database.Config) error {
	var dsn string
	if config.DBName != "" {
		dsn = config.DBName
	} else {
		// 使用内存数据库
		// Use in-memory database
		dsn = ":memory:"
	}

	var err error
	s.db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	return nil
}

// Close 关闭 SQLite 数据库连接
// Close closes the SQLite database connection.
func (s *SQLiteDB) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

// DB 返回 GORM 的数据库实例
// DB returns the GORM DB instance.
func (s *SQLiteDB) DB() *gorm.DB {
	return s.db
}
