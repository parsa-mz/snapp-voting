package controllers

import (
	"SnappVotingBack/app/services"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
)

type FileController struct {
}

func (FileController) Serve(ctx *gin.Context) {
	object, err := services.MinioClient.GetObject(context.Background(),
		"reportage-snapp",
		ctx.Param("file_name"), minio.GetObjectOptions{})
	if err != nil {
		print(err.Error())
		ctx.JSON(404, gin.H{
			"message": "error",
		})
		return
	}
	_, err = object.Stat()
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "error",
		})
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, object)
	if err != nil {
		print(err.Error())
		ctx.JSON(500, gin.H{
			"message": "error",
		})
		return
	} // Error handling elided for brevity.

	ctx.Data(200, http.DetectContentType(buf.Bytes()), buf.Bytes())
}
