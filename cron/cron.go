package cron

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Cron struct {
	cron   *cron.Cron
	tgBot  *bot.Bot
	logger *zerolog.Logger
}

func NewCron(l *zerolog.Logger, b *bot.Bot) *Cron {
	logger := l.With().Str("type", "bot").Logger()

	c := &Cron{
		cron:   cron.New(),
		tgBot:  b,
		logger: &logger,
	}

	if _, err := c.cron.AddFunc("@hourly", diskUsageJob(&logger, b)); err != nil {
		logger.Err(err).Msg("Failed to schedule job")
	}

	return c
}

func (c *Cron) Start() {
	c.logger.Info().Msg("Starting cron")
	c.cron.Start()
}

func diskUsageJob(l *zerolog.Logger, b *bot.Bot) func() {
	return func() {
		logger := l.With().Str("type", "diskUsageJob").Logger()

		cmd := exec.Command("sh", "-c", "df -h /host")
		var out, stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.Err(err).Str("stderr", stderr.String()).Msg("Failed to execute df")
			return
		}

		for _, line := range strings.Split(out.String(), "\n") {
			if strings.HasSuffix(line, "/host") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					usage := fields[4]
					msg := fmt.Sprintf("Disk usage: %s", usage)
					logger.Info().Msg(msg)
					b.SendMessage(msg)
					return
				}
			}
		}

		logger.Warn().Msg("Could not parse /host usage from df output")
	}
}
