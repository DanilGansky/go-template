package api

import (
	"context"
	"net/http"
	"time"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/auth"
	"github.com/littlefut/go-template/internal/user"

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

func NewAuthController(authSvc auth.Service, userSvc user.Service, opts *Options) *AuthController {
	c := &AuthController{
		Engine:  opts.Router,
		authSvc: authSvc,
		userSvc: userSvc,
		log:     opts.Log,
		timeout: opts.Timeout,
	}

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
		ctx.JSON(http.StatusBadRequest, errors.New(errors.ValidationError, err))
		return
	}

	token, err := c.authSvc.Login(newCtx, &dto)
	if err != nil {
		var status int

		switch errors.Cause(err) {
		case errors.NotFoundError:
			status = http.StatusNotFound
		case errors.ValidationError:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.log.Errorf("login error: %s with status: %d", err.Error(), status)
		ctx.JSON(status, err)
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
		ctx.JSON(http.StatusBadRequest, errors.New(errors.ValidationError, err))
		return
	}

	err := c.userSvc.Register(newCtx, &dto)
	if err != nil {
		var status int

		switch errors.Cause(err) {
		case errors.NotFoundError:
			status = http.StatusNotFound
		case errors.ValidationError:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.log.Errorf("register error: %s with status: %d", err.Error(), status)
		ctx.JSON(status, err)
		return
	}

	ctx.Status(http.StatusCreated)
	c.log.Infof("user with username: '%s' created", dto.Username)
}
