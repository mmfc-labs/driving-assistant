package apiserver

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"github.com/mmfc-labs/driving-assistant/version"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mmfc-labs/driving-assistant/docs"
	log "github.com/sirupsen/logrus"
)

type APIServer struct {
	router   *gin.Engine
	srv      *http.Server
	validate *validator.Validate
	lbs      *lbs.LBS
}

func NewAPIServer(opt Options) *APIServer {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	router := gin.Default()

	srv := &http.Server{
		Addr:    opt.Addr,
		Handler: router,
	}

	apiServer := &APIServer{
		router:   router,
		srv:      srv,
		validate: validate,
	}

	err := config.LoadConfig(opt.ConfigPath, func(setting config.Setting, probe probe.Probe) {
		apiServer.lbs = lbs.NewLBS(setting, probe)
		log.WithField("setting", setting).WithField("probe", probe).Info("重新加载配置成功")
	})

	if err != nil {
		panic(err)
	}

	apiServer.registerAPI()
	return apiServer
}

func (s *APIServer) registerAPI() {
	s.router.LoadHTMLGlob("./templates/*")
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "coord.html", gin.H{
			"title": "coord",
		})
	})

	s.router.GET("/api/healthz", func(c *gin.Context) {
		c.String(200, "I'm fine")
	})
	s.router.GET("/api/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": version.Version, "gitRevision": version.GitRevision})
	})
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.Use(HandleCors).GET("/api/route", s.route)
	s.router.Use(HandleCors).GET("/api/probes", s.probes)
}

func (s *APIServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("api server forced to shutdown:", err)
	}
	log.Info("api server exit ")
}

func (s *APIServer) Run() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("api server close")
				return
			}
			log.Fatal("api server listen: ", err)
		}
	}()
}

func Result(httpStatus int, data interface{}, errorMsg string, c *gin.Context) {
	c.JSON(httpStatus, Response{
		data,
		errorMsg,
	})
	return
}

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}
