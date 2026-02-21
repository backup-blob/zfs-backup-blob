package config

type CompressConfig struct {
	T        string `yaml:"type" validate:"required"`
	Remote_  string `yaml:"remote" validate:"required"`
	Level    int    `yaml:"level" validate:"min=1,max=22"`
}

func NewCompressConfig() ConfigStage {
	return &CompressConfig{
		Level: 3, // default level: good balance of speed/compression
	}
}

func (c *CompressConfig) Type() ConfigType {
	return Middleware
}

func (c *CompressConfig) Remote() string {
	return c.Remote_
}
