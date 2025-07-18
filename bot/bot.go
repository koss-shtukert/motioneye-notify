package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/koss-shtukert/motioneye-notify/pkg"
	"github.com/rs/zerolog"
)

type Bot struct {
	tgBot  *tgbotapi.BotAPI
	chatId int64
	logger *zerolog.Logger
	upload *pkg.Upload
}

func CreateBot(key, cid string, l *zerolog.Logger, u *pkg.Upload) (*Bot, error) {
	logger := l.With().Str("type", "bot").Logger()

	tgBot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	chatId, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id: %w", err)
	}

	return &Bot{
		tgBot:  tgBot,
		chatId: chatId,
		logger: &logger,
		upload: u,
	}, nil
}

func (bot *Bot) SendMessage(m string) {
	msg := tgbotapi.NewMessage(bot.chatId, m)
	if _, err := bot.tgBot.Send(msg); err != nil {
		bot.logger.Err(err).Msg("Failed to send message")
	}
}

func (bot *Bot) SendPhoto(caption, url string) {
	data, err := bot.upload.UploadFromUrl(url)
	if err != nil {
		bot.logger.Err(err).Msg("Failed to upload photo from URL")
		return
	}

	photo := tgbotapi.NewPhoto(bot.chatId, tgbotapi.FileBytes{
		Name:  caption,
		Bytes: data,
	})
	photo.Caption = caption

	if _, err := bot.tgBot.Send(photo); err != nil {
		bot.logger.Err(err).Msg("Failed to send photo")
	}
}
