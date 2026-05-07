
package postgres

import (
	"go-clean-architecture/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB 鏍规嵁閰嶇疆鏂囦欢鍒濆鍖栧苟杩斿洖 PostgreSQL 杩炴帴瀹炰緥
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	if cfg.Postgres.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	}
	if cfg.Postgres.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	}

	return db, nil
}
