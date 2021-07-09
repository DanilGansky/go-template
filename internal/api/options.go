package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Log     *logrus.Logger
	Timeout time.Duration
	Router  *gin.Engine
}
