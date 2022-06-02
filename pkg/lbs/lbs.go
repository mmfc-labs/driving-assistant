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

// LBS 根据路面距离计算需要避让的探头
type LBS struct {
	client  drive.Client
	setting config.Setting
	probe   probe.Probe
}

// NewLBS
func NewLBS(setting config.Setting, probe probe.Probe) *LBS {
	c := &LBS{
		client:  tencent.NewClient(config.TencentKey),
		setting: setting,
		probe:   probe,
	}
	return c
}

//Route 根据直线距离计算需要避让的探头
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
func (c *LBS) Route(from, to drive.Coord) (avoidAreas [][]drive.Coord, debug *apis.Debug, err error) {
	probesMap := make(map[drive.Coord]struct{}, 0)
	debug = &apis.Debug{}
	debug.RouteLogs = make([]*apis.DebugLog, 0)
	count := 0
Again:
	count++
	debugLog := &apis.DebugLog{}

	if count > c.setting.MaxRoute {
		return avoidAreas, debug, apis.RouteOutOfRange
	}

	probes := make([]drive.Coord, 0, len(probesMap))
	isAgain := false

	for coord, _ := range probesMap {
		probes = append(probes, coord)
	}

	route, err := c.client.GetRoutes(from, to, probes, c.setting.AvoidAreaOffset)
	if err != nil {
		return avoidAreas, debug, err
	}
	routePoints := route[0].Points

	// Debug Info
	debugLog.Debug(count, routePoints, probes)

	debugCurToNextToProbe := make([]string, 0)
	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]
		probePoints := c.probe.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint drive.Coord) {
			if c.isAvoid(cur, next, probePoint) {
				curToNextToProbeInfo := fmt.Sprintf("needAvoid,fmt(cur;next;probe): %s", drive.FmtCoord(cur, next, probePoint))
				log.Info(curToNextToProbeInfo)
				debugCurToNextToProbe = append(debugCurToNextToProbe, curToNextToProbeInfo)
				probesMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}

	// Debug Info
	debug.Debug(debugLog, debugCurToNextToProbe, count, probesMap)

	if isAgain {
		log.Info("again,count:", count)
		goto Again
	}
	if len(probesMap) > c.setting.MaxAvoid {
		return avoidAreas, debug, apis.AvoidOutOfRange
	}

	avoidAreas = make([][]drive.Coord, 0)
	for a := range probesMap {
		avoidAreas = append(avoidAreas, drive.ConvCoordToAvoidArea(a, config.AvoidAreaOffset))
	}
	log.WithField("avoidPoints", probesMap).WithField("count", count).Info("根据直线距离计算需要避让的探头")
	return avoidAreas, debug, nil
}

func (c *LBS) isAvoid(cur drive.Coord, next drive.Coord, probePoint drive.Coord) bool {
	_, b1, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
	_, b2, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon})
	_, b3, _ := geodist.VincentyDistance(geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
	gap := b1 - (b2 + b3 - (float64(c.setting.Offset) / 1000))
	log.Infof("calculate: b1:%f, b2:%f, b3:%f, offset:%f, gap:%f", b1, b2, b3, float64(c.setting.Offset)/1000, gap)
	return gap > 0
}

func (c *LBS) Probes(cur drive.Coord, near float64) ([]drive.Coord, error) {
	if near == 0 {
		return c.probe.Points, nil
	}
	return c.probe.Near(cur, near), nil
}
