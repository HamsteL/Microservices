package routers

import (
	"github.com/gin-gonic/gin"
)

func HomeGET(c *gin.Context) {
	c.String(200, "Home page")
}
