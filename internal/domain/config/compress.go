package config

// DefaultCompressLevel is the default compression level (balance of speed/ratio).
const DefaultCompressLevel = 3

// CompressConfig holds compression settings.
type CompressConfig struct {
	T       string `yaml:"type" validate:"required"`
	Remote_ string `yaml:"remote" validate:"required"`
	Level   int    `yaml:"level" validate:"min=1,max=22"`
}

// NewCompressConfig creates a new compress config with default level.
func NewCompressConfig() ConfigStage {
	return &CompressConfig{
		Level: DefaultCompressLevel,
	}
}

// Type returns the config type for this stage.
func (c *CompressConfig) Type() ConfigType {
	return Middleware
}

// Remote returns the remote stage name.
func (c *CompressConfig) Remote() string {
	return c.Remote_
}
