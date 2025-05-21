package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"APP_ENV"`
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	ServerHost    string `mapstructure:"SERVER_HOST"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	MotioneyeHost string `mapstructure:"MOTIONEYE_HOST"`
	MotioneyePort string `mapstructure:"MOTIONEYE_PORT"`
	TgBotApiKey   string `mapstructure:"TGBOT_API_KEY"`
	TgBotChatId   string `mapstructure:"TGBOT_CHAT_ID"`
}

func Load(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("viper couldn't read in the config file. %v", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("viper could not unmarshal the configuration. %v", err)
	}

	if config.Environment == "" {
		return nil, errors.New("APP_ENV variable is empty")
	}

	if config.LogLevel == "" {
		return nil, errors.New("LOG_LEVEL variable is empty")
	}

	if config.ServerHost == "" {
		return nil, errors.New("SERVER_HOST variable is empty")
	}

	if config.ServerPort == "" {
		return nil, errors.New("SERVER_PORT variable is empty")
	}

	if config.MotioneyeHost == "" {
		return nil, errors.New("MOTIONEYE_HOST variable is empty")
	}

	if config.MotioneyePort == "" {
		return nil, errors.New("MOTIONEYE_PORT variable is empty")
	}

	if config.TgBotApiKey == "" {
		return nil, errors.New("TGBOT_API_KEY variable is empty")
	}

	if config.TgBotChatId == "" {
		return nil, errors.New("TGBOT_CHAT_ID variable is empty")
	}

	return
}
