package config

type S3Config struct {
	t              string `yaml:"type" validate:"required"`
	Bucket         string `yaml:"bucket"`
	Region         string `yaml:"region"`
	UsePathStyle   bool   `yaml:"usePathStyle"`
	BaseEndpoint   string `yaml:"baseEndpoint"`
	Prefix         string `yaml:"prefix"`
	AccessKey      string `yaml:"accessKey"`
	AccessSecret   string `yaml:"accessSecret"`
	MaxRetries     int    `yaml:"maxRetries"`
	UploadPartSize int    `yaml:"uploadPartSize"`
}

func NewS3Config() ConfigStage {
	return &S3Config{}
}

func (s *S3Config) Type() ConfigType {
	return Sink
}

func (s *S3Config) Remote() string {
	return ""
}
