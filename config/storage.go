package config

type Storage struct {
	Tmp   string     `mapstructure:"tmp" json:"tmp" yaml:"tmp"` // 临时存储
	Way   storageWay `mapstructure:"tmp" json:"way" yaml:"way"`
	Local local      `mapstructure:"lcaol" json:"local" yaml:"local"`
}
type storageWay uint32

const (
	LocalStorage storageWay = iota // 本地存储
)

type local struct {
	Dir string `mapstructure:"dir" json:"dir" yaml:"dir"`
}
