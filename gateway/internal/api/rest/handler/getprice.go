package get

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetPrice(c *gin.Context) {
	res, err := h.pricingSvc.GetPrice(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}
