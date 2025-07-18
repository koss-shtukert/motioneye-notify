package callback

import (
	"github.com/koss-shtukert/motioneye-notify/bot"
	"github.com/koss-shtukert/motioneye-notify/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Index(e *echo.Echo, cc *config.Config, b *bot.Bot) {
	e.GET("/", func(c echo.Context) (err error) {
		cameraName := c.QueryParam("camera_name")
		cameraId := c.QueryParam("camera_id")

		b.SendPhoto(cameraName, getCameraSnapUrl(cameraId, cc))

		return c.String(http.StatusOK, "Ok")
	})
}

func getCameraSnapUrl(ci string, c *config.Config) string {
	Url := "http://" + c.MotioneyeHost + ":" + c.MotioneyePort
	snapUrl := Url + "/picture/" + ci + "/current/"

	return snapUrl
}
