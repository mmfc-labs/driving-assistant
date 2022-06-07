package apis

import (
	"fmt"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	log "github.com/sirupsen/logrus"
)

type Debug struct {
	RouteCount int         `json:"route_count"` // 路线规划次数
	ProbeCount int         `json:"probe_count"` // 需要避让的探头数量
	RouteLogs  []*DebugLog `json:"route_logs"`  // 路线规划日志
}

func (debug *Debug) Debug(debugLog *DebugLog, debugCurToNextToProbe []string, count int, probesMap map[drive.Coord]struct{}) {
	debugLog.CurToNextToProbe = debugCurToNextToProbe
	debug.RouteLogs = append(debug.RouteLogs, debugLog)
	debug.ProbeCount = len(probesMap)
	debug.RouteCount = count
}

type DebugLog struct {
	RouteInfo        string   `json:"route_info"`           //本次路线信息
	RouteProbeInfo   string   `json:"route_probe_info"`     //本次路线信息传入的避让探头, 第一次传入的为空
	CurToNextToProbe []string `json:"cur_to_next_to_probe"` //本次计算后需要避让的探头，格式为：cur;next;probe (A1;A2;探头)
}

func (d *DebugLog) Debug(count int, routePoints []drive.Coord, probes []drive.Coord) {
	routeInfo := fmt.Sprintf("第%d次路线规划,坐标串:%s", count, drive.FmtCoord(routePoints...))
	d.RouteInfo = routeInfo
	d.RouteProbeInfo = fmt.Sprintf("第%d次路线规划,传入的探头:%s", count, drive.FmtCoord(probes...))
	log.Info(d.RouteInfo)
	log.Info(d.RouteProbeInfo)
}
