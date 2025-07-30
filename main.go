package main

import (
	"log"

	"github.com/koss-shtukert/motioneye-notify/api"
	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/koss-shtukert/motioneye-notify/logger"
	"github.com/koss-shtukert/motioneye-notify/pkg"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	logr, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal("Logger error: ", err)
	}

	uploader := pkg.CreateUploader(&logr)

	tgBot, err := bot.CreateBot(cfg.TgBotApiKey, cfg.TgBotChatId, &logr, uploader)
	if err != nil {
		log.Fatal("Telegram bot error: ", err)
	}

	logr.Info().Str("type", "core").Msg("Starting service")

	server := api.CreateServer(&logr, cfg, tgBot)

	log.Fatal(server.Start(":" + cfg.ServerPort))
}
