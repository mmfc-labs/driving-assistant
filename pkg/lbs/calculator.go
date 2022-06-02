package lbs

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/stream"
)

// Calculator 根据路面距离计算需要避让的探头
type Calculator struct {
	client  drive.Client
	setting config.Setting
	probe   probe.Probe
}

// NewCalculator
func NewCalculator(setting config.Setting, probe probe.Probe) *Calculator {
	c := &Calculator{
		client:  tencent.NewClient(config.TencentKey),
		setting: setting,
		probe:   probe,
	}
	return c
}

//AvoidProbe 根据直线距离计算需要避让的探头
//首先获取到路线A坐标串
//A1 -> A2 -> A3 -> A4 -> AN
//
// 需要转弯在十字路口都会有坐标串
//
//A1 -> A2   直线距离 B1
//A1 -> 探头1 直线距离 B2
//探头1 -> A2 直线距离 B3
//
//B1 >= B2+B3 即路过探头
func (c *Calculator) AvoidProbe(from, to drive.Coord) (probesMap map[drive.Coord]struct{}, debug *apis.DebugResp, err error) {
	logEntry := log.WithField("route", "AvoidProbe")

	probesMap = make(map[drive.Coord]struct{}, 0)
	debug = &apis.DebugResp{}
	debug.Routes = make([]apis.DebugRouteResp, 0)
	count := 0
Again:
	debugRoute := apis.DebugRouteResp{}

	count++
	if count > c.setting.MaxRoute {
		return probesMap, debug, fmt.Errorf("超出最大路线规划次数:%d", c.setting.MaxRoute)
	}

	probes := make([]drive.Coord, 0, len(probesMap))
	isAgain := false

	for coord, _ := range probesMap {
		probes = append(probes, coord)
	}

	route, err := c.client.GetRoutes(from, to, probes, c.setting.AvoidAreaOffset)
	if err != nil {
		return probesMap, debug, err
	}
	routePoints := route[0].Points
	routeInfo := fmt.Sprintf("第%d次路线:%s", count, drive.FmtCoord(routePoints...))
	logEntry.Info(routeInfo)
	debugRoute.RouteInfo = routeInfo
	debugRoute.RouteProbeInfo = fmt.Sprintf("第%d次路线,传入的探头:%s", count, drive.FmtCoord(probes...))

	debugCurToNextToProbe := make([]string, 0)
	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]
		probePoints := c.probe.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint drive.Coord) {
			_, b1, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
			_, b2, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon})
			_, b3, _ := geodist.VincentyDistance(geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
			gap := b1 - (b2 + b3 - (float64(c.setting.Offset) / 1000))
			logEntry.Infof("calculate: b1:%f, b2:%f, b3:%f, offset:%f, gap:%f", b1, b2, b3, float64(c.setting.Offset)/1000, gap)
			if gap > 0 {
				curToNextToProbeInfo := fmt.Sprintf("needAvoid,fmt(cur;next;probe): %s", drive.FmtCoord(cur, next, probePoint))
				logEntry.Info(curToNextToProbeInfo)
				debugCurToNextToProbe = append(debugCurToNextToProbe, curToNextToProbeInfo)
				probesMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}

	debugRoute.CurToNextToProbe = debugCurToNextToProbe
	debug.Routes = append(debug.Routes, debugRoute)
	debug.ProbeCount = len(probesMap)
	debug.RouteCount = count
	if isAgain {
		logEntry.Info("again,count:", count)
		goto Again
	}
	if len(probesMap) > c.setting.MaxAvoid {
		return probesMap, debug, fmt.Errorf("超出最大避让点：CurAvoid:%d MaxAvoid:%d", len(probesMap), c.setting.MaxAvoid)
	}

	logEntry.Info("执行次数：", count)
	return probesMap, debug, nil
}

func (c *Calculator) Probes(cur drive.Coord, near float64) ([]drive.Coord, error) {
	if near == 0 {
		return c.probe.Points, nil
	}
	return c.probe.Near(cur, near), nil
}
