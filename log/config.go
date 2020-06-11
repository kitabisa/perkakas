package log

import "github.com/spf13/viper"

type Config struct {
	Log Log
}
type Log struct {
	Level	string `mapstructure:"level"`
}

func newConfig() (conf *Config, err error) {
	err = viper.Unmarshal(&conf)
	if err != nil {
		return
	}
	return
}