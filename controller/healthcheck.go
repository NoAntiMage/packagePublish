package controller

import (
	"PackageServer/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	service.NewTodoService()
	c.String(http.StatusOK, "ok")
}
