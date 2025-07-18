package api

import (
	"github.com/koss-shtukert/motioneye-notify/api/rest/callback"
	"github.com/koss-shtukert/motioneye-notify/api/rest/test"
	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func CreateServer(l *zerolog.Logger, c *config.Config, b *bot.Bot) *echo.Echo {
	logger := l.With().Str("type", "api").Logger()

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log := logger.Debug()

			if v.Error != nil {
				log = logger.Error()
			}

			log.
				Str("host", v.Host).
				Str("uri", v.URI).
				Str("method", v.Method).
				Int("status", v.Status).
				Any("headers", v.Headers).
				Str("remote_ip", v.RemoteIP).
				Str("request_id", v.RequestID)

			if v.Error == nil {
				log.Msg("request")
			} else {
				log.Msg(v.Error.Error())
			}

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"GET"},
	}))

	e.HideBanner = true

	callback.REST(e, c, b)
	test.REST(e, c, b)

	return e
}
