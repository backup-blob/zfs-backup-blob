package repository

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
)

type configRepo struct {
	ConfigDriver config.ConfigDriver
}

func (c *configRepo) GetMiddlewares() []domain.Middleware {
	return c.ConfigDriver.GetMiddlewares()
}

func NewConfig(configDriver config.ConfigDriver) config.ConfigRepo {
	return &configRepo{ConfigDriver: configDriver}
}

func (c *configRepo) GetConfig() *config.Config {
	return c.ConfigDriver.GetConfig()
}
