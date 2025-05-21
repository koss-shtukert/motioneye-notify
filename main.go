package main

import (
	"github.com/koss-shtukert/motioneye-notify/pkg"
	"log"

	"github.com/koss-shtukert/motioneye-notify/api"
	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/koss-shtukert/motioneye-notify/logger"
)

func main() {
	c, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.New(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	u := pkg.CreateUploader(&l)

	b, err := bot.CreateBot(c.TgBotApiKey, c.TgBotChatId, &l, u)
	if err != nil {
		log.Fatal(err)
	}

	l.Info().Str("type", "core").Msg("Starting service")

	s := api.CreateServer(&l, c, b)

	log.Fatal(s.Start(":" + c.ServerPort))
}
