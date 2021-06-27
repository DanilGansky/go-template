package api

import (
	"context"
	"net/http"
	"time"

	"github.com/littlefut/go-template/internal/auth"
	"github.com/littlefut/go-template/internal/user"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	*gin.Engine

	authSvc auth.Service
	userSvc user.Service

	log     *logrus.Logger
	timeout time.Duration
}

func NewAuthController(authSvc auth.Service, userSvc user.Service, log *logrus.Logger, timeout time.Duration, r *gin.Engine) *AuthController {
	c := &AuthController{r, authSvc, userSvc, log, timeout}

	group := c.Group("/auth")
	group.POST("/login", c.Login)
	group.POST("/register", c.Register)
	return c
}

func (c *AuthController) Login(ctx *gin.Context) {
	newCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var dto auth.LoginDTO
	if err := ctx.Bind(&dto); err != nil {
		c.log.Errorf("login error: %s", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	token, err := c.authSvc.Login(newCtx, &dto)
	if err != nil {
		var status int

		switch errors.Cause(err) {
		case user.ErrNotFound:
			status = http.StatusNotFound
		case user.ErrValidation:
			status = http.StatusBadRequest
		case auth.ErrInvalidPassword:
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
		}

		c.log.Errorf("login error: %s with status: %d", err.Error(), status)
		ctx.Status(status)
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (c *AuthController) Register(ctx *gin.Context) {
	newCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var dto user.RegisterDTO
	if err := ctx.Bind(&dto); err != nil {
		c.log.Errorf("register error: %s", err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := c.userSvc.Register(newCtx, &dto)
	if err != nil {
		var status int

		switch errors.Cause(err) {
		case user.ErrNotFound:
			status = http.StatusNotFound
		case user.ErrValidation:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.log.Errorf("register error: %s with status: %d", err.Error(), status)
		ctx.Status(status)
		return
	}

	ctx.Status(http.StatusCreated)
	c.log.Infof("user with username: '%s' created", dto.Username)
}
