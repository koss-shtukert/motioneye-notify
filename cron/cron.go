package cron

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"strconv"
)

type Cron struct {
	cron   *cron.Cron
	tgBot  *bot.Bot
	logger *zerolog.Logger
	config *config.Config
}

func NewCron(l *zerolog.Logger, cfg *config.Config, b *bot.Bot) *Cron {
	logger := l.With().Str("type", "cron").Logger()

	c := &Cron{
		cron:   cron.New(),
		tgBot:  b,
		logger: &logger,
		config: cfg,
	}

	if _, err := c.cron.AddFunc(cfg.CronDiskUsageJobInterval, diskUsageJob(&logger, cfg, b)); err != nil {
		logger.Err(err).Msg("Failed to schedule job")
	}

	return c
}

func (c *Cron) Start() {
	c.logger.Info().Msg("Starting cron")
	c.cron.Start()
}

func diskUsageJob(l *zerolog.Logger, c *config.Config, b *bot.Bot) func() {
	return func() {
		logger := l.With().Str("type", "diskUsageJob").Logger()

		path := "/host" + c.CronDiskUsageJobPath
		cmd := exec.Command("sh", "-c", "df -h "+path)

		var out, stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			logger.Err(err).Str("stderr", stderr.String()).Msg("Failed to execute df")
			return
		}

		for _, line := range strings.Split(out.String(), "\n") {
			if strings.HasSuffix(line, path) {
				fields := strings.Fields(line)
				if len(fields) < 5 {
					logger.Warn().Str("line", line).Msg("Not enough fields to parse")
					continue
				}

				used := fields[2]
				avail := fields[3]
				usedPercent := fields[4]

				percentStr := strings.TrimSuffix(usedPercent, "%")
				percent, err := strconv.Atoi(percentStr)
				if err != nil {
					logger.Warn().Str("value", usedPercent).Msg("Could not parse Use% value")
					return
				}

				msg := fmt.Sprintf("Disk\nUsed: %s\nAvail: %s\nUse%%: %d%%", used, avail, percent)
				logger.Info().
					Str("path", path).
					Str("used", used).
					Str("avail", avail).
					Int("use_percent", percent).
					Msg("Parsed disk usage")
				b.SendMessage(msg)
				return
			}
		}

		logger.Warn().Msgf("Could not parse disk usage for path %s from df output", path)
	}
}
