package bot

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/koss-shtukert/motioneye-notify/pkg"
	"github.com/rs/zerolog"
	"strconv"
)

type Bot struct {
	tgBot  *tgbotapi.BotAPI
	chatId int64
	logger *zerolog.Logger
	upload *pkg.Upload
}

func CreateBot(key string, cid string, l *zerolog.Logger, u *pkg.Upload) (*Bot, error) {
	logger := l.With().Str("type", "bot").Logger()
	var err error

	tgBot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		err = errors.New("Error creating bot" + err.Error())
	}

	return &Bot{
		tgBot:  tgBot,
		chatId: func() int64 { id, _ := strconv.ParseInt(cid, 10, 64); return id }(),
		logger: &logger,
		upload: u,
	}, err
}

func (bot *Bot) SendMessage(m string) {
	msg := tgbotapi.NewMessage(bot.chatId, m)

	if _, err := bot.tgBot.Send(msg); err != nil {
		bot.logger.Error().Err(err).Msg("Failed to send message: " + err.Error())
	}
}

func (bot *Bot) SendPhoto(c string, u string) {
	pb, err := bot.upload.UploadFromUrl(u)

	if err != nil {
		bot.logger.Error().Err(err).Msg("Failed to send photo: " + err.Error())
	}

	if pb != nil {
		photo := tgbotapi.NewPhoto(bot.chatId, tgbotapi.FileBytes{
			Name:  c,
			Bytes: pb,
		})

		photo.Caption = c

		if _, err := bot.tgBot.Send(photo); err != nil {
			bot.logger.Error().Err(err).Msg("Failed to send photo: " + err.Error())
		}
	}
}
