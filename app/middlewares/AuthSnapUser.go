package middlewares

import (
	"SnappVotingBack/app/models"
	"SnappVotingBack/app/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthSnappUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		snappUser := new(models.SnappUser)
		var err error
		snappUser.SnappId = ctx.Param("snapp_id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, serializers.Base{
				Code:    serializers.InvalidInput,
				Message: "snapp_id must be int64",
			})
			return
		}
		exists, err := snappUser.GetUser()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, serializers.Base{
				Code:    serializers.InternalError,
				Message: "internal server error, Connection ERR R-M-ASU-IE",
			})
		}
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusNotFound, serializers.Base{
				Code:    serializers.SnappIdDoesNotExists,
				Message: "snapp_id does not exists",
			})
		}
		ctx.Set("snappUser_id", snappUser.Id)
		ctx.Next()
	}
}
