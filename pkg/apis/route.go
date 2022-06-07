package apis

import "github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"

type RouteReq struct {
	FromLat float64 `form:"from_lat" json:"from_lat" validate:"required"`
	FromLon float64 `form:"from_lon" json:"from_lon" validate:"required"`
	ToLat   float64 `form:"to_lat" json:"to_lat" validate:"required"`
	ToLon   float64 `form:"to_lon" json:"to_lon" validate:"required" label:"to_lon"`
}

type RouteResp struct {
	AvoidAreas  [][]drive.Coord `json:"avoid_areas"`  // 需要避让的区域
	AvoidProbes []drive.Coord   `yaml:"avoid_probes"` // 需要避让的探头
	Debug       *Debug          `json:"debug"`        //debug信息
}
