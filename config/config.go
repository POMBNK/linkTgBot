package config

import (
	"sync"

	"github.com/POMBNK/linktgBot/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token    string `yml:"token"`
	DbPath   string `yml:"db_path"`
	Host     string `yml:"host"`
	LogLevel string `yml:"log_level"`
}

var cfg *Config
var once sync.Once

func GetCfg() *Config {
	once.Do(func() {
		logger := logging.GetLogger("info")
		logger.Info("Считываем конфиг...")
		cfg = &Config{}
		if err := cleanenv.ReadConfig("config.yml", cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return cfg
}
