package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // CGO-free SQLite driver

	"github.com/432539/gpt2api/internal/config"
)

// NewSQLite 打开(或创建)SQLite 数据库并返回连接池。
// 自动创建父目录;应用 WAL 模式和合理的 PRAGMA 提升并发性能。
func NewSQLite(cfg config.SQLiteConfig) (*sqlx.DB, error) {
	path := cfg.Path
	if path == "" {
		path = "data/gpt2api.db"
	}

	// 确保父目录存在
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create sqlite dir: %w", err)
		}
	}

	// DSN 附加 SQLite 参数:
	//   _journal=WAL       — WAL 模式,读写并发友好
	//   _timeout=5000      — 忙等超时 5s,避免 SQLITE_BUSY 立即报错
	//   _foreign_keys=on   — 开启外键约束
	//   _cache=shared      — 共享缓存(同进程多连接)
	dsn := fmt.Sprintf(
		"file:%s?_journal=WAL&_timeout=5000&_foreign_keys=on&_cache=shared",
		path,
	)

	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite 不支持真正的并发写;限制连接池到合理大小,
	// 设置 MaxOpenConns=1 可彻底避免 SQLITE_BUSY,
	// WAL 模式下读可以并发,写串行由 Go mutex 保护。
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}
