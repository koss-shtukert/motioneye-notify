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

		b.SendPhoto(cameraName, getCameraSnapUrl(cameraName, cc))

		return c.String(http.StatusOK, "Ok")
	})
}

func getCameraSnapUrl(cn string, c *config.Config) string {
	Url := "http://" + c.MotioneyeHost + ":" + c.MotioneyePort
	snapUrl := ""

	switch cn {
	case "AlxFront":
		snapUrl = Url + "/picture/1/current/"
		break
	case "KossFront":
		snapUrl = Url + "/picture/2/current/"
		break
	case "KossBack":
		snapUrl = Url + "/picture/3/current/"
		break
	}

	return snapUrl
}
