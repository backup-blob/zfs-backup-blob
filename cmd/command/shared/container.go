package shared

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	"github.com/backup-blob/zfs-backup-blob/pkg/format"
	"github.com/golobby/container/v3"
	"io"
	"os"
	"os/user"
	"strconv"
	"time"
)

const ConfigFileName = ".bbackup.yaml"

type nower struct {
}

func (n *nower) Now() time.Time {
	timeFlag := os.Getenv("BB_FLAG_TIME")
	if timeFlag != "" {
		i, err := strconv.ParseInt(timeFlag, 10, 64)
		if err != nil {
			panic(err)
		}

		return time.Unix(i, 0)
	}

	return time.Now()
}

func getSizer() func(*int64) string {
	fixSize := os.Getenv("BB_FLAG_FIX_SIZE")
	if fixSize != "" {
		return func(i *int64) string {
			return "1MB"
		}
	}

	return format.Size
}

func loadConfig(configPath string) (io.Reader, error) {
	if configPath == "" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("failed to get current user workingdir %w", err)
		}

		configPath = currentUser.HomeDir + "/" + ConfigFileName
	}

	file, err := os.Open(configPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read config file %w", err)
	}

	return file, nil
}

func LoadDeps(configPath, logLevel string) container.Container {
	c := container.New() //nolint:varnamelen // no need

	// driver
	container.MustSingleton(c, func() domain.LogDriver {
		return driver.NewLog(os.Stdout, domain.StringToLevel(logLevel))
	})
	container.MustSingleton(c, func() domain.SnapshotNamestrategy {
		return driver.NewDefaultNamer(&nower{})
	})
	container.MustSingleton(c, driver.NewRender)
	container.MustSingletonLazy(c, func(logger domain.LogDriver) (config.ConfigDriver, error) {
		configReader, err := loadConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config %w", err)
		}

		return driver.NewConfigDriver(&config.LoadParams{
			ConfigReader: configReader,
			StageMapping: map[string]func() config.ConfigStage{
				"s3":       config.NewS3Config,
				"crypt":    config.NewCryptConfig,
				"zfs":      config.NewZfsConfig,
				"throttle": config.NewThrottleConfig,
			},
			StorageDriverFunc: func(s *config.S3Config) (domain.StorageDriver, error) {
				return driver.NewS3StorageFromConfig(s, logger)
			},
			ZfsDriverFunc: func(cZfs *config.ZfsConfig) domain.ZfsDriver {
				return driver.NewZfsFromConfig(cZfs, logger)
			},
			ToMiddleware: func(c config.ConfigStage) domain.Middleware {
				switch v := c.(type) {
				case *config.ThrottleConfig:
					return driver.NewThrottle(v)
				case *config.CryptConfig:
					return driver.NewCrypt(v)
				default:
					return nil
				}
			},
		})
	})

	// repos
	container.MustSingletonLazy(c, repository.NewLog)
	container.MustSingletonLazy(c, func(renderDriver domain.RenderDriver) domain.RenderRepository {
		return repository.NewRender(renderDriver, getSizer())
	})
	container.MustSingletonLazy(c, func(confDriver config.ConfigDriver, namer domain.SnapshotNamestrategy) domain.SnapshotRepository {
		return repository.NewSnapshot(confDriver.GetZfsDriver(), namer)
	})
	container.MustSingletonLazy(c, func(confDriver config.ConfigDriver) domain.VolumeRepository {
		return repository.NewVolume(confDriver.GetZfsDriver())
	})
	container.MustSingletonLazy(c, func(confDriver config.ConfigDriver) domain.BackupRepository {
		return repository.NewBackup(confDriver.GetZfsDriver(), confDriver.GetStorageDriver())
	})
	container.MustSingletonLazy(c, func(confDriver config.ConfigDriver) domain.BackupStateRepo {
		return repository.NewBackupStateRepo(confDriver.GetStorageDriver())
	})
	container.MustSingletonLazy(c, repository.NewConfig)

	// usecase
	container.MustSingletonLazy(c, usecase.NewSnapshot)
	container.MustSingletonLazy(c, usecase.NewBackup)
	container.MustSingletonLazy(c, usecase.NewBackupList)
	container.MustSingletonLazy(c, usecase.NewVolumeUsecase)
	container.MustSingletonLazy(c, usecase.NewBackupSync)
	container.MustSingletonLazy(c, usecase.NewTrimUseCase)

	return c
}
