package apiserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// route
// @Tags driving
// @Summary 路线规划，获取需要避让的区域
// @accept application/json
// @Produce application/json
// @Param data query RouteRequest true "请求"
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
	avoidPoints, err := s.calculator.AvoidProbe(from, to)
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

// probes
// @Tags driving
// @Summary 获取探头
// @accept application/json
// @Produce application/json
// @Param data query ProbeRequest true "请求"
// @success 200 {object} Response{data=ProbeResponse} "返回结果"
// @Router /api/probes [get]
func (s *APIServer) probes(c *gin.Context) {
	var (
		req ProbeRequest
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}
	if err := s.validate.Struct(req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}
	probes, err := s.calculator.Probes(drive.Coord{Lat: req.Lat, Lon: req.Lon}, req.Near)
	if err != nil {
		Result(http.StatusInternalServerError, nil, err.Error(), c)
		return
	}

	Result(http.StatusOK, ProbeResponse{Probes: probes}, "", c)
}

type ProbeRequest struct {
	Lat  float64 `form:"lat" json:"lat"`   // 当前位置经度
	Lon  float64 `form:"lon" json:"lon"`   // 当前位置纬度
	Near float64 `form:"near" json:"near"` // 获取附近多少公里的探头，0为获取所有探头
}

type ProbeResponse struct {
	Probes []drive.Coord `json:"probes"`
}
