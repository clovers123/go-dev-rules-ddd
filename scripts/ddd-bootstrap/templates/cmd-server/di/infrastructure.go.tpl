package di

import (
	"time"

	"git.yugeeker.com/SHARED/go-lazy/cache"
	"git.yugeeker.com/SHARED/go-lazy/config"
	"git.yugeeker.com/SHARED/go-lazy/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"{{.ModulePath}}/app/acl/adapter/repository/postgres/repository"
	"{{.ModulePath}}/configs"
)

var InfrastructureModule = fx.Options(
	fx.Provide(
		configs.LoadConfig,
		log.ProvideLogger,
		func(cfg *configs.Config) (redis.UniversalClient, error) {
			return redis.NewClient(&redis.Options{Addr: "localhost:6379"}), nil
		},
		func(cfg *configs.Config) (*gorm.DB, error) {
			dsn, err := cfg.PostgresDSN("default")
			if err != nil {
				return nil, err
			}
			return gorm.Open(postgres.Open(dsn), &gorm.Config{})
		},
		func(cfg *config.BasicConfiguration) *cache.LruCache {
			return cache.NewLruCache(1024, time.Hour)
		},
	),
	fx.Invoke(func(db *gorm.DB) {
		repository.SetDefault(db)
	}),
)
