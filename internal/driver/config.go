package driver

import (
	"errors"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"io"
)

type ConfigDriver struct {
	zfsDriver     domain.ZfsDriver
	storageDriver domain.StorageDriver
	middlewares   []domain.Middleware
	config        *config.Config
}

func NewConfigDriver(params *config.LoadParams) (config.ConfigDriver, error) {
	conf := &ConfigDriver{}

	err := conf.Load(params)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *ConfigDriver) unmarshal(r io.Reader, mapping config.StageMap) (*config.Config, error) {
	var configS config.Config

	yamlFile, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &configS)
	if err != nil {
		return nil, fmt.Errorf("unmarshalFile failed %w", err)
	}

	err = c.unmarshalStages(&configS, mapping)
	if err != nil {
		return nil, err
	}

	return &configS, nil
}

func (c *ConfigDriver) unmarshalStages(conf *config.Config, mapping map[string]func() config.ConfigStage) error {
	conf.Stages = make(map[string]config.ConfigStage)

	for key, stage := range conf.RawStages { //nolint:gocritic // copy wont impact speed much
		var stageS struct {
			Type string `yaml:"type"`
		}

		err := stage.Decode(&stageS)
		if err != nil {
			return fmt.Errorf("failed to read stage: %s", key)
		}

		loader, ok := mapping[stageS.Type]
		if !ok {
			return fmt.Errorf("config type for stage: %s not found", key)
		}

		entity := loader()

		errA := stage.Decode(entity)
		if errA != nil {
			return fmt.Errorf("config type %s failed", key)
		}

		conf.Stages[key] = entity
	}

	return nil
}

func (c *ConfigDriver) orderStages(conf *config.Config) ([]config.ConfigStage, error) {
	var (
		list   []config.ConfigStage
		source string
	)

	graph := NewGraph()

	for key, item := range conf.Stages {
		remote := item.Remote()

		if remote != "" {
			_, exists := conf.Stages[remote]
			if !exists {
				return nil, fmt.Errorf("remote not found: %s", remote)
			}

			graph.AddEdge(key, remote)
		}

		if item.Type() == config.Source {
			if source != "" {
				return nil, fmt.Errorf("there can only be one source")
			}

			source = key
		}
	}

	if source == "" {
		return nil, fmt.Errorf("source not found")
	}

	if graph.hasCycle() {
		return nil, errors.New("stages cannot have cycles")
	}

	for _, key := range graph.DFS(source) {
		item, ok := conf.Stages[key]
		if !ok {
			return nil, fmt.Errorf("stages not found %s", key)
		}

		list = append(list, item)
	}

	return list, nil
}

func (c *ConfigDriver) GetMiddlewares() []domain.Middleware {
	return c.middlewares
}

func (c *ConfigDriver) GetStorageDriver() domain.StorageDriver {
	return c.storageDriver
}

func (c *ConfigDriver) GetZfsDriver() domain.ZfsDriver {
	return c.zfsDriver
}

func (c *ConfigDriver) GetConfig() *config.Config {
	return c.config
}

func (c *ConfigDriver) Load(params *config.LoadParams) error {
	conf, err := c.unmarshal(params.ConfigReader, params.StageMapping)
	if err != nil {
		return fmt.Errorf("configRepo.loadFileConfig failed %w", err)
	}

	c.config = conf

	err = c.validateStruct()
	if err != nil {
		return err
	}

	orderedStages, err := c.orderStages(conf)
	if err != nil {
		return fmt.Errorf("configRepo.orderStages failed %w", err)
	}

	err = c.loadStages(params, orderedStages)
	if err != nil {
		return err
	}

	err = c.validateDrivers()
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigDriver) validateDrivers() error {
	if c.zfsDriver == nil {
		return fmt.Errorf("zfs missing")
	}

	if c.storageDriver == nil {
		return fmt.Errorf("storage missing")
	}

	return nil
}

func (c *ConfigDriver) validateStruct() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c.config)

	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (c *ConfigDriver) loadStages(params *config.LoadParams, stages []config.ConfigStage) error {
	for index, stage := range stages {
		switch stage.Type() {
		case config.Source:
			c.fillZfs(stage, params)
		case config.Sink:
			if err := c.fillS3(stage, params); err != nil {
				return err
			}
		case config.Middleware:
			res := params.ToMiddleware(stage)
			if res == nil {
				return fmt.Errorf("invalid middleware on stage index %d", index)
			}

			c.middlewares = append(c.middlewares, res)
		}
	}

	return nil
}

func (c *ConfigDriver) fillS3(stage config.ConfigStage, params *config.LoadParams) error {
	newS3, ok := stage.(*config.S3Config)
	if ok {
		storage, err := params.StorageDriverFunc(newS3)
		if err != nil {
			return err
		}

		c.storageDriver = storage
	}

	return nil
}

func (c *ConfigDriver) fillZfs(stage config.ConfigStage, params *config.LoadParams) {
	newZfs, ok := stage.(*config.ZfsConfig)
	if ok {
		c.zfsDriver = params.ZfsDriverFunc(newZfs)
	}
}
