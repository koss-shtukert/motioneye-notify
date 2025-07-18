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
				if len(fields) >= 5 {
					used := fields[2]
					avail := fields[3]
					usageStr := fields[4]

					percent := 0
					if _, err := fmt.Sscanf(usageStr, "%d%%", &percent); err != nil {
						l.Err(err).Str("raw", usageStr).Msg("Failed to parse usage percentage")
						return
					}

					var msg string

					msg += "💾 Disk Usage\n\n"
					msg += "+--------+---------+--------+\n"
					msg += "| Used   | Avail   | Use%%  |\n"
					msg += "+--------+---------+--------+\n"
					msg += fmt.Sprintf("| %-6s | %-7s | %-5s |\n", used, avail, usageStr)
					msg += "+--------+---------+--------+"

					if percent >= 90 {
						msg = "🚨 *High Disk Usage Alert!*\n\n" + msg
					}

					b.SendMessage(msg)
					return
				}
			}
		}

		logger.Warn().Msgf("Could not parse disk usage from df output on path %s", path)
	}
}
