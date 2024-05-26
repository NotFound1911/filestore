package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Configuration 配置文件中所有字段对应的结构体
type Configuration struct {
	Log     Log     `mapstructure:"log" json:"log" yaml:"log"`
	Mq      Mq      `mapstructure:"mq" json:"mq" yaml:"mq"`
	Storage Storage `mapstructure:"storage" json:"storage" yaml:"storage"`
}

const confFilePath = "conf/config.yaml"

func NewConfig(path string) *Configuration {
	conf := &Configuration{}
	if path == "" {
		initConfig(confFilePath, conf)
	} else {
		initConfig(path, conf)
	}
	return conf
}
func initConfig(path string, conf *Configuration) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		//lgLogger.Logger.Error("read config failed: ", zap.String("err", err.Error()))
		fmt.Println("read config failed: ", zap.String("err", err.Error()))
		panic(err)
	}

	if err := v.Unmarshal(&conf); err != nil {
		//lgLogger.Logger.Error("config parse failed: ", zap.String("err", err.Error()))
		fmt.Println("config parse failed: ", zap.String("err", err.Error()))
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		//lgLogger.Logger.Info("", zap.String("config file changed:", in.Name))
		fmt.Println("", zap.String("config file changed:", in.Name))
		defer func() {
			if err := recover(); err != nil {
				//lgLogger.Logger.Error("config file changed err:", zap.Any("err", err))
				fmt.Println("config file changed err:", zap.Any("err", err))
			}
		}()
		if err := v.Unmarshal(&conf); err != nil {
			//lgLogger.Logger.Error("config parse failed: ", zap.String("err", err.Error()))
			fmt.Println("config parse failed: ", zap.String("err", err.Error()))
		}
	})
}
