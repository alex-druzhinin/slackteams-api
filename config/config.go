package config

import (
	"fmt"

	"github.com/crgimenes/goconfig"
)

func LoadConfig() (config Config, err error) {
	goconfig.PrefixEnv = "ST_API"

	if err := goconfig.Parse(&config); err != nil {
		return config, err
	}

	return config, validateEnv(config)
}

func validateEnv(config Config) error {
	envs := []EnvType{EnvTypeDevelopment, EnvTypeStage, EnvTypeProduction}

	for _, env := range envs {
		if env == config.Env {
			return nil
		}
	}

	return fmt.Errorf("Unknown Env: '%s'", config.Env)
}

type Config struct {
	Env        EnvType       `cfgRequired:"true"`
	Addr       string        `cfgDefault:":7010"`
	MongoDB    MongoDBConfig `cfgRequired:"true"`
	Amqp       AmqpConfig
	BotUser    User `cfgRequired:"true"`
	MeteorUser User `cfgRequired:"true"`
}

type User struct {
	Token string
	Name  string
}

type MongoDBConfig struct {
	URI string `cfgRequired:"true"`
}

type AmqpConfig struct {
	URI string `cfgRequired:"true"`
}

type EnvType string

const (
	EnvTypeDevelopment EnvType = "development"
	EnvTypeStage       EnvType = "stage"
	EnvTypeProduction  EnvType = "production"
)
