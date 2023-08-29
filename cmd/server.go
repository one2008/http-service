package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"http-service/cmd/log"
)

var (
	Logger log.Logger = nil
	Conf   *Config    = nil
)

type Server struct {
	logger log.Logger
	engine *gin.Engine
	addr   string
}

type ginLogger struct {
	logger log.Logger
}

func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

func (g ginLogger) Write(p []byte) (n int, err error) {
	g.logger.Info(string(p))
	return len(p), nil
}

func NewServer(conf *Config, gorm *gorm.DB, logger log.Logger) (*Server, error) {
	Logger = logger
	Conf = conf
	router := gin.New()
	router.Use(gin.LoggerWithWriter(ginLogger{logger: logger}))
	router.Use(gin.Recovery())
	router.Use(NoCache())

	return &Server{
		logger: logger,
		engine: router,
		addr:   conf.Server.Addr,
	}, nil
}

func (s *Server) Router() {
	v1 := s.engine.Group(InternalApi_prefix_V1)
	{
		v1.Use(AuthMiddleware())
		for _, v := range ApiPathList {
			if v.Method == http.MethodGet {
				v1.GET(v.Path, v.Handler)
			} else if v.Method == http.MethodPost {
				v1.POST(v.Path, v.Handler)
			}
		}
		/*
			// 内部服务调用接口
			// 内部接口，apikey验证
			v1.POST("/example", func(ctx *gin.Context) {})

			// 前端接口
			// 服务器unix时间戳
			v1.GET("/serverUnixTimestamp", func(ctx *gin.Context) {
				resp := common.BizResponse{}
				common.SetupSuccess(&resp, map[string]interface{}{"serverTime": time.Now().UTC().Unix()})
				ctx.JSON(http.StatusOK, resp)
			})

		*/
	}
}

func (s *Server) Run() error {
	InitApiPath(s)

	s.Router()

	return s.engine.Run(s.addr)
}
