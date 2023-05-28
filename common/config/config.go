package config

import (
	"main/utils/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"false"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"tcp"`
		BindIp string `yaml:"bind_ip" env-default:"0.0.0.0"`
		Port   string `yaml:"port" env-default:"3000"`
	} `yaml:"listen"`
	MongoDB struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"27017"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"  env-default:"test"`
		AuthDB   string `yaml:"authDB"`
	} `yaml:"mongodb"`
	Redis struct {
		Host    string `yaml:"host" env-default:"localhost"`
		Port    string `yaml:"port" env-default:"6379"`
		Passwod string `yaml:"password" env-default:""`
		DB      int    `yaml:"DB" env-default:"0"`
	} `yaml:"redis"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read appliaction configuration")

		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
