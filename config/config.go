package config

import (
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

func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	required := map[string]string{
		"APP_ENV":        cfg.Environment,
		"LOG_LEVEL":      cfg.LogLevel,
		"SERVER_HOST":    cfg.ServerHost,
		"SERVER_PORT":    cfg.ServerPort,
		"MOTIONEYE_HOST": cfg.MotioneyeHost,
		"MOTIONEYE_PORT": cfg.MotioneyePort,
		"TGBOT_API_KEY":  cfg.TgBotApiKey,
		"TGBOT_CHAT_ID":  cfg.TgBotChatId,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("required variable %s is empty", key)
		}
	}

	return &cfg, nil
}
