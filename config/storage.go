package config

type Storage struct {
	Tmp   string     `mapstructure:"tmp" json:"tmp" yaml:"tmp"` // 临时存储
	Way   storageWay `mapstructure:"way" json:"way" yaml:"way"`
	Local local      `mapstructure:"local" json:"local" yaml:"local"`
	Minio minio      `mapstructure:"minio" json:"minio" yaml:"minio"`
}
type storageWay uint32

const (
	LocalStorage storageWay = iota // 本地存储
	MinioStorage                   // minio存储
)

type local struct {
	Dir string `mapstructure:"dir" json:"dir" yaml:"dir"`
}

type minio struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secret_access_key" yaml:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl" json:"use_ssl" yaml:"use_ssl"`
}
