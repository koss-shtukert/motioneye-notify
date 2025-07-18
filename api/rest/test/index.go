package test

import (
	"net/http"

	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/labstack/echo/v4"
)

func Index(e *echo.Echo, cfg *config.Config, b *bot.Bot) {
	e.GET("/test", handleTest(cfg, b))
}

func handleTest(cfg *config.Config, b *bot.Bot) echo.HandlerFunc {
	return func(c echo.Context) error {
		cameraName := c.QueryParam("camera_name")
		cameraID := c.QueryParam("camera_id")

		b.SendPhoto(cameraName, getCameraSnapURL(cameraID, cfg))

		return c.String(http.StatusOK, "Ok")
	}
}

func getCameraSnapURL(cameraID string, cfg *config.Config) string {
	return "http://" + cfg.MotioneyeHost + ":" + cfg.MotioneyePort + "/picture/" + cameraID + "/current/"
}
