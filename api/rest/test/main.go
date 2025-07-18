package test

import (
	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/labstack/echo/v4"
)

func REST(e *echo.Echo, c *config.Config, b *bot.Bot) {
	Index(e, c, b)
}
