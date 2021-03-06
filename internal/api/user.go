package api

import (
	"context"
	"net/http"
	"time"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	*gin.Engine

	userSvc user.Service
	log     *logrus.Logger
	timeout time.Duration
}

func NewUserController(userSvc user.Service, opts *Options, middlewares ...gin.HandlerFunc) *UserController {
	c := &UserController{
		Engine:  opts.Router,
		userSvc: userSvc,
		log:     opts.Log,
		timeout: opts.Timeout,
	}

	group := c.Group("/users")
	group.Use(middlewares...)
	group.PATCH("/me", c.UpdateUser)
	group.DELETE("/me", c.DeleteUser)
	return c
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	newCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var dto user.UpdateDTO
	if err := ctx.Bind(&dto); err != nil {
		c.log.Errorf("update user error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, errors.New(errors.ValidationError, err))
		return
	}

	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrInvalidContext)
		return
	}

	userID, ok := id.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrInvalidContext)
		return
	}

	err := c.userSvc.SetUsername(newCtx, userID, &dto)
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

		c.log.Errorf("update user error: %s with status: %d", err.Error(), status)
		ctx.JSON(status, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	newCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	id, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrInvalidContext)
		return
	}

	userID, ok := id.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrInvalidContext)
		return
	}

	err := c.userSvc.Delete(newCtx, userID)
	if err != nil {
		var status int

		switch errors.Cause(err) {
		case errors.NotFoundError:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}

		c.log.Errorf("delete user error: %s with status: %d", err.Error(), status)
		ctx.JSON(status, err)
		return
	}

	c.log.Infof("user deleted with id: %d", id)
	ctx.Status(http.StatusNoContent)
}
