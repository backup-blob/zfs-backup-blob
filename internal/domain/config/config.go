package config

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"gopkg.in/yaml.v3"
	"io"
)

type ConfigType int

const (
	Sink ConfigType = iota
	Middleware
	Source
)

type StageMap = map[string]func() ConfigStage

type ConfigDriver interface {
	GetZfsDriver() domain.ZfsDriver
	GetStorageDriver() domain.StorageDriver
	GetMiddlewares() []domain.Middleware
	GetConfig() *Config
}

type ConfigRepo interface {
	GetConfig() *Config
	GetMiddlewares() []domain.Middleware
}

type LoadParams struct {
	ConfigReader      io.Reader
	StageMapping      StageMap
	ZfsDriverFunc     func(c *ZfsConfig) domain.ZfsDriver
	StorageDriverFunc func(s *S3Config) (domain.StorageDriver, error)
	ToMiddleware      func(c ConfigStage) domain.Middleware
}

type ConfigStage interface {
	Remote() string
	Type() ConfigType
}

type Config struct {
	RemoteTrimPolicy RemoteTrimPolicy       `yaml:"remote_trim_policy"`
	LocalTrimPolicy  LocalTrimPolicy        `yaml:"local_trim_policy"`
	Stages           map[string]ConfigStage `yaml:"-" validate:"required,dive"`
	RawStages        map[string]yaml.Node   `yaml:"stages"`
}
