package config

type ThrottleConfig struct {
	t          string `yaml:"type" validate:"required"`
	Remote_    string `yaml:"remote" validate:"required"`
	WriteSpeed int64  `yaml:"writeSpeed" validate:"required"` //TODO: rename to writeSpeedByteSec
	ReadSpeed  int64  `yaml:"readSpeed" validate:"required"`  //TODO: rename to readSpeedByteSec
}

func NewThrottleConfig() ConfigStage {
	return &ThrottleConfig{}
}

func (t *ThrottleConfig) Remote() string {
	return t.Remote_
}

func (t *ThrottleConfig) Type() ConfigType {
	return Middleware
}
