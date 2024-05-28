package config

type Log struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Dir        string `mapstructure:"dir" json:"dir" yaml:"dir"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine   bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // day
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
	EnableFile bool   `mapstructure:"enable_file" json:"enable_file" yaml:"enable_file"`
}
