package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Interval time.Duration
	SSH      *SSH
	Metrics  *Metrics
}

type SSH struct {
	Host     string
	Username string
	Password string
	Port     uint16
}

type Metrics struct {
	Port uint16
	Path string
}

var config *Config

func LoadConfig() error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()

	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	if config == nil {
		LoadConfig()
	}

	return config
}
