package config

import (
	"log"

	"github.com/spf13/viper"

	"github.com/kiritoxkiriko/comical-tool/pkg/constant"
)

var (
	Conf Config
)

type Config struct {
	Server   Server   `mapstructure:"server" json:"server"`
	Database Database `mapstructure:"database" json:"database"`
	Log      Log      `mapstructure:"log" json:"log"`
	Redis    Redis    `mapstructure:"redis" json:"redis"`
}

type Server struct {
	Port int `mapstructure:"port" json:"port"`
}

type Database struct {
	DSN string `mapstructure:"dsn" json:"dsn"`
}

type Log struct {
	Level string `mapstructure:"level" json:"level"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

func InitConfig() {
	viper.AddConfigPath(constant.ConfigPath)
	viper.AddConfigPath(constant.ConfigPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config file failed: %v", err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("unmarshal config failed: %v", err)
	}
}
