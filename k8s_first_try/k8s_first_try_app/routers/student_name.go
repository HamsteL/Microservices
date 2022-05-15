package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func HelloStudent(c *gin.Context) {
	var studentName string = c.Param("name")
	c.String(200, fmt.Sprintf("Hello %s!", studentName))
}
