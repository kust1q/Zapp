package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type handlerError struct {
	Message string `json:"message"`
}

func newErrorReponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, handlerError{message})
}
