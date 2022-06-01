package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"github.com/mmfc-labs/driving-assistant/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mmfc-labs/driving-assistant/docs"
	log "github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type APIServer struct {
	router     *gin.Engine
	srv        *http.Server
	validate   *validator.Validate
	calculator *lbs.Calculator
}

func NewAPIServer(opt Options) *APIServer {
	apiServer := &APIServer{
		validate: validator.New(),
	}

	router := gin.Default()
	router.GET("/api/healthz", func(c *gin.Context) {
		c.String(200, "I'm fine")
	})
	router.GET("/api/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": version.Version, "gitRevision": version.GitRevision})
	})
	router.Use(HandleCors).GET("/api/route", apiServer.route)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    opt.Addr,
		Handler: router,
	}

	apiServer.router = router
	apiServer.srv = srv

	err := config.LoadConfig(opt.ConfigPath, func(setting config.Setting, probe probe.Probe) {
		apiServer.calculator = lbs.NewCalculator(setting, probe)
		log.Info("重新加载配置成功")
	})

	if err != nil {
		panic(err)
	}

	return apiServer
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

// route
// @Tags driving
// @Summary 路线规划，获取需要避让的区域
// @accept application/json
// @Produce application/json
// @Param data query RouteRequest true "RouteRequest"
// @success 200 {object} Response{data=RouteResponse} "返回结果"
// @Router /api/route [get]
func (s *APIServer) route(c *gin.Context) {
	var (
		req RouteRequest
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}
	if err := s.validate.Struct(req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}

	// 起点，终点
	from, to := drive.Coord{Lat: req.FromLat, Lon: req.FromLon}, drive.Coord{Lat: req.ToLat, Lon: req.ToLon}

	//根据直线距离计算需要避让的探头
	avoidPoints, err := s.calculator.AvoidProbeByLine(from, to)
	if err != nil {
		Result(http.StatusInternalServerError, nil, err.Error(), c)
		return
	}
	log.Info("根据直线距离计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}

	avoidArea := make([][]drive.Coord, 0)
	for a, _ := range avoidPoints {
		avoidArea = append(avoidArea, drive.ConvCoordToAvoidArea(a, config.AvoidAreaOffset))
	}

	Result(http.StatusOK, RouteResponse{AvoidAreas: avoidArea}, "", c)
}

type RouteRequest struct {
	FromLat float64 `form:"from_lat" json:"from_lat" validate:"required"`
	FromLon float64 `form:"from_lon" json:"from_lon" validate:"required"`
	ToLat   float64 `form:"to_lat" json:"to_lat" validate:"required"`
	ToLon   float64 `form:"to_lon" json:"to_lon" validate:"required" label:"to_lon"`
}

type RouteResponse struct {
	AvoidAreas [][]drive.Coord `json:"avoid_areas"`
}

func Result(httpStatus int, data interface{}, errorMsg string, c *gin.Context) {
	c.JSON(httpStatus, Response{
		data,
		errorMsg,
	})
	return
}

type Response struct {
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"error_msg"`
}
