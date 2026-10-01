package read

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

type Health interface {
	Get(c *gin.Context)
}
