package config

type Mq struct {
	Addr           []string `mapstructure:"addr" json:"addr" yaml:"addr"`
	NetDiaTimeout  int      `mapstructure:"net_dia_timeout" json:"net_dia_timeout" yaml:"net_dia_timeout"`
	NetReadTimeout int      `mapstructure:"net_read_timeout" json:"net_read_timeout" yaml:"net_read_timeout"`
	PReqiredAcks   int      `mapstructure:"p_required_acks" json:"p_required_acks" yaml:"p_required_acks"`
	PTimeout       int      `mapstructure:"p_timeout" json:"p_timeout" yaml:"p_timeout"`
	CStrategy      string   `mapstructure:"c_strategy" json:"c_strategy" yaml:"c_strategy"`
	CMaxWaitTime   int      `mapstructure:"c_max_wait_time" json:"c_max_wait_time" yaml:"c_max_wait_time"`
	CFetchMin      int      `mapstructure:"c_fetch_min" json:"c_fetch_min" yaml:"c_fetch_min"`
	CFetchDefault  int      `mapstructure:"c_fetch_default" json:"c_fetch_default" yaml:"c_fetch_default"`
	PReturnSuccess bool     `mapstructure:"p_return_success" json:"p_return_success" yaml:"p_return_success"`
	PReturnErr     bool     `mapstructure:"p_return_err" json:"p_return_err" yaml:"p_return_err"`
	Enable         bool     `mapstructure:"enable" json:"enable" yaml:"enable"`
}
