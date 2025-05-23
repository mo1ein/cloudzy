package get

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetWeather(c *gin.Context) {
	res, err := h.weatherSvc.GetWeather(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}
