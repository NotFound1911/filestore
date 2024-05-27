package config

type Etcd struct {
	Endpoints []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
}
