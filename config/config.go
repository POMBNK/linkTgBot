package config

import (
	"sync"

	"github.com/POMBNK/linktgBot/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token  string `yml:"token"`
	DbPath string `yml:"db_path"`
	Host   string `yml:"host"`
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
			logger.Info("Инф из YML")
			logger.Info(help)
			logger.Fatal(err)
		}
		// if err := cleanenv.ReadEnv(cfg); err != nil {
		// 	help, _ := cleanenv.GetDescription(cfg, nil)
		// 	logger.Info("Инф из ENV")
		// 	logger.Info(help)
		// 	logger.Fatal(err)
		// }
	})
	return cfg
}
