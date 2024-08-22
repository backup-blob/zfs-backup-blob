package config

type ZfsConfig struct {
	t       string `yaml:"type" validate:"required"`
	Remote_ string `yaml:"remote" validate:"required"`
	ZfsPath string `yaml:"zfsPath"`
}

func NewZfsConfig() ConfigStage {
	return &ZfsConfig{}
}

func (z *ZfsConfig) Remote() string {
	return z.Remote_
}

func (z *ZfsConfig) Type() ConfigType {
	return Source
}
