package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ImageProxy struct {
}

func (*ImageProxy) Handler(ctx *gin.Context) {
	l := struct {
		Url string `json:"url" form:"url"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", l.Url, nil)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	_, _ = io.Copy(ctx.Writer, resp.Body)
	resp.Body.Close()
	return
}

func NewImageProxy() *ImageProxy {
	return &ImageProxy{}
}
