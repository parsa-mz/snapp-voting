package controllers

import (
	"SnappVotingBack/app/models"
	"SnappVotingBack/app/serializers"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type User struct {
}

func (User) Register(ctx *gin.Context) {
	var request serializers.UserRequest
	if ctx.BindJSON(&request) != nil {
		ctx.JSON(400, serializers.Base{Message: serializers.InvalidInput})
		return
	}
	user := models.User{
		Email: request.Email,
	}
	err := user.SetPassword(request.Password)
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(500, serializers.Base{Message: serializers.InternalError})
		return
	}
	err = user.Create()
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(500, serializers.Base{Message: serializers.InternalError})
		return
	}
	ctx.JSON(200, serializers.Base{Message: serializers.Success})
}

func (User) Login(ctx *gin.Context) {
	var request serializers.UserRequest
	if ctx.BindJSON(&request) != nil {
		ctx.JSON(400, serializers.Base{Message: serializers.InvalidInput})
		return
	}
	user := models.User{
		Email: request.Email,
	}
	err := user.Get()
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(404, serializers.Base{Message: serializers.NotFound})
		return
	}
	Access := user.CheckPassword(request.Password)
	if !Access {
		ctx.JSON(403, serializers.Base{Message: serializers.WrongPassword})
		return
	}
	auth, err := user.Auth()
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(500, serializers.Base{Message: serializers.InternalError})
		return
	}
	ctx.JSON(200, serializers.UserJWT{
		Access: auth,
	})
}

func (User) Reset(ctx *gin.Context) {
	var request serializers.UserRequest
	if ctx.BindJSON(&request) != nil {
		ctx.JSON(400, serializers.Base{Message: serializers.InvalidInput})
		return
	}
	var user models.User
	user.Id = ctx.GetInt64("user_id")

	err := user.SetPassword(request.Password)
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(500, serializers.Base{Message: serializers.InternalError})
		return
	}
	err = user.UpdatePassword()
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(500, serializers.Base{Message: serializers.InternalError})
		return
	}
	ctx.JSON(200, serializers.Base{Message: serializers.Success})
}
