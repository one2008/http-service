package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	InternalApi_prefix_V1 = "v1"
)

type ApiPath struct {
	Path     string
	Method   string
	Isfilter bool
	Handler  gin.HandlerFunc
}

var ApiPathList = []ApiPath{}

func InitApiPath(s *Server) {
	ApiPathList = []ApiPath{
		{
			"/example",
			http.MethodPost,
			true,
			func(ctx *gin.Context) {

			},
		},

		{
			"/serverUnixTimestamp",
			http.MethodGet,
			false,
			func(ctx *gin.Context) {
				resp := BizResponse{}
				SetupSuccess(&resp, map[string]interface{}{"serverTime": time.Now().UTC().Unix()})
				ctx.JSON(http.StatusOK, resp)
			},
		},
	}
}
