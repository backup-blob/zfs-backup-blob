package config

type CryptConfig struct {
	t        string `yaml:"type" validate:"required"`
	Remote_  string `yaml:"remote" validate:"required"`
	Password string `yaml:"password" validate:"required"` //nolint:gosec // configuration field for encryption password
}

func NewCryptConfig() ConfigStage {
	return &CryptConfig{}
}

func (c *CryptConfig) Type() ConfigType {
	return Middleware
}

func (c *CryptConfig) Remote() string {
	return c.Remote_
}
