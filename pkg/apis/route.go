package apis

import "github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"

type RouteReq struct {
	FromLat float64 `form:"from_lat" json:"from_lat" validate:"required"`
	FromLon float64 `form:"from_lon" json:"from_lon" validate:"required"`
	ToLat   float64 `form:"to_lat" json:"to_lat" validate:"required"`
	ToLon   float64 `form:"to_lon" json:"to_lon" validate:"required" label:"to_lon"`
}

type RouteResp struct {
	AvoidAreas [][]drive.Coord `json:"avoid_areas"`
	Debug      *DebugResp      `json:"debug"`
}

type DebugResp struct {
	RouteCount int             `json:"route_count"` // 路线规划次数
	ProbeCount int             `json:"probe_count"` // 需要避让的探头数量
	RouteLogs  []DebugLogsResp `json:"route_logs"`  // 路线规划日志
}

type DebugLogsResp struct {
	RouteInfo        string   `json:"route_info"`           //本次路线信息
	RouteProbeInfo   string   `json:"route_probe_info"`     //本次路线信息传入的避让探头, 第一次传入的为空
	CurToNextToProbe []string `json:"cur_to_next_to_probe"` //本次计算后需要避让的探头，格式为：cur;next;probe (A1;A2;探头)
}
