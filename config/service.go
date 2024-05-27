package config

type Service struct {
	Apigw       Apigw       `mapstructure:"apigw" json:"apigw" yaml:"apigw"`
	Account     Account     `mapstructure:"account" json:"account" yaml:"account"`
	FileManager FileManager `mapstructure:"file_manager" json:"file_manager" yaml:"file_manager"`
	Upload      Upload      `mapstructure:"upload" json:"upload" yaml:"upload"`
	Transfer    Transfer    `mapstructure:"transfer" json:"transfer" yaml:"transfer"`
}

type Apigw struct {
	Name string    `mapstructure:"name" json:"name" yaml:"name"`
	Http webServer `mapstructure:"http" json:"http" yaml:"http"`
}

type Account struct {
	Name string     `mapstructure:"name" json:"name" yaml:"name"`
	Grpc grpcServer `mapstructure:"grpc" json:"grpc" yaml:"grpc"`
}

type FileManager struct {
	Name string     `mapstructure:"name" json:"name" yaml:"name"`
	Grpc grpcServer `mapstructure:"grpc" json:"grpc" yaml:"grpc"`
}

type Upload struct {
	Name string    `mapstructure:"name" json:"name" yaml:"name"`
	Http webServer `mapstructure:"http" json:"http" yaml:"http"`
}

type Transfer struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type webServer struct {
	Mode string   `mapstructure:"mode" json:"mode" yaml:"mode"`
	Addr []string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type grpcServer struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}
