package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mmfc-labs/driving-assistant/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"github.com/mmfc-labs/driving-assistant/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type APIServer struct {
	router   *gin.Engine
	srv      *http.Server
	validate *validator.Validate
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
	router.Use(HandleCors).GET("/api/avoids", apiServer.avoids)

	srv := &http.Server{
		Addr:    opt.Addr,
		Handler: router,
	}

	apiServer.router = router
	apiServer.srv = srv
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

func (s *APIServer) avoids(c *gin.Context) {
	var (
		req AvoidsRequest
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
	calculator := lbs.NewCalculator(tencent.NewClient(config.TencentKey), config.Offset, config.AvoidAreaOffset)

	//根据直线距离计算需要避让的探头
	avoidPoints, err := calculator.AvoidProbeByLine(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println("根据直线距离计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}

	avoidArea := make([][]drive.Coord, 0)
	for a, _ := range avoidPoints {
		avoidArea = append(avoidArea, drive.ConvCoordToAvoidArea(a, config.AvoidAreaOffset))
	}
	Result(http.StatusOK, avoidArea, "", c)
}

type AvoidsRequest struct {
	FromLat float64 `form:"from_lat" validate:"required"`
	FromLon float64 `form:"from_lon" validate:"required"`
	ToLat   float64 `form:"to_lat" validate:"required"`
	ToLon   float64 `form:"to_lon" validate:"required" label:"to_lon"`
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

type AvoidsResponse struct {
	AvoidArea [][]drive.Coord
}
