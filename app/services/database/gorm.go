package database

import (
	"context"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// reader is a readonly database connection pool.
	reader *gorm.DB
	// writer is a full-access database connection pool.
	writer *gorm.DB
)

// GetReader fetches a gorm readonly connection, allocating one if needed.
func GetReader(ctx context.Context) (*gorm.DB, error) {
	if reader == nil {
		r, err := newGormConnection(ctx, config.MustConfig()["database.reader"])
		if err != nil {
			return nil, err
		}
		reader = r
	}
	return reader, nil
}

// GetWriter fetches a gorm connection, allocating one if needed.
func GetWriter(ctx context.Context) (*gorm.DB, error) {
	if writer == nil {
		w, err := newGormConnection(ctx, config.MustConfig()["database.writer"])
		if err != nil {
			return nil, err
		}
		writer = w
	}
	return writer, nil
}

// newGormConnection instantiates a new database connection.
func newGormConnection(ctx context.Context, dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
