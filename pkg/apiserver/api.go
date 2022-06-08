package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"net/http"
)

// route
// @Tags driving
// @Summary 路线规划，获取需要避让的区域
// @accept application/json
// @Produce application/json
// @Param data query apis.RouteReq true "请求"
// @success 200 {object} Response{data=apis.RouteResp} "返回结果"
// @Router /api/route [get]
func (s *APIServer) route(c *gin.Context) {
	var (
		req apis.RouteReq
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
	from, to := geo.Coord{Lat: req.FromLat, Lon: req.FromLon}, geo.Coord{Lat: req.ToLat, Lon: req.ToLon}
	//根据直线距离计算需要避让的探头
	avoidAreas, avoidProbes, debug, err := s.lbs.Route(from, to)
	if err != nil {
		Result(http.StatusInternalServerError, apis.RouteResp{Debug: debug}, err.Error(), c)
		return
	}
	Result(http.StatusOK, apis.RouteResp{AvoidAreas: avoidAreas, AvoidProbes: avoidProbes, Debug: debug}, "", c)
}

// probes
// @Tags driving
// @Summary 获取探头
// @accept application/json
// @Produce application/json
// @Param data query apis.ProbeReq true "请求"
// @success 200 {object} Response{data=apis.ProbeResp} "返回结果"
// @Router /api/probes [get]
func (s *APIServer) probes(c *gin.Context) {
	var (
		req apis.ProbeReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}
	if err := s.validate.Struct(req); err != nil {
		Result(http.StatusBadRequest, nil, err.Error(), c)
		return
	}
	probes, err := s.lbs.Probes(geo.Coord{Lat: req.Lat, Lon: req.Lon}, req.Near)
	if err != nil {
		Result(http.StatusInternalServerError, nil, err.Error(), c)
		return
	}

	Result(http.StatusOK, apis.ProbeResp{Probes: probes}, "", c)
}
