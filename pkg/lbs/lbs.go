package lbs

import (
	"fmt"
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/stream"
	"math"
)

// LBS 根据路面距离计算需要避让的探头
type LBS struct {
	client       drive.Client
	setting      config.Setting
	probeManager probe.ProbeManager
}

// NewLBS
func NewLBS(setting config.Setting, probeManager probe.ProbeManager) *LBS {
	c := &LBS{
		client:       tencent.NewClient(setting.LBSKey),
		setting:      setting,
		probeManager: probeManager,
	}
	return c
}

//Route 根据直线距离计算需要避让的探头
//首先获取到路线A坐标串,不传入避让区
//A1 -> A2 -> A3 -> A4 -> AN. 需要转弯的都会在十字路口都会有坐标串
//
//
//条件一:
//
//A1 -> A2   直线距离 B1
//A1 -> 探头1 直线距离 B2
//探头1 -> A2 直线距离 B3
//B1 >= B2+B3-offset 即路过探头
//
// 条件二:
//
//A1->A2             朝向toward1 （角度值）
//探头->探头朝向坐标    朝向toward2 （角度值）
//toward1 与 toward2 差值小于 TowardRange 则视为同一个方向
//
//满足条件一与条件二则需要避让该探头
func (c *LBS) Route(from, to drive.Coord) (avoidAreas [][]drive.Coord, avoidProbes []drive.Coord, debug *apis.Debug, err error) {
	probesMap := make(map[drive.Coord]struct{}, 0)
	debug = &apis.Debug{}
	debug.RouteLogs = make([]*apis.DebugLog, 0)
	count := 0
Again:
	count++
	debugLog := &apis.DebugLog{}

	if count > c.setting.MaxRoute {
		return avoidAreas, avoidProbes, debug, apis.RouteOutOfRange
	}

	probes := make([]drive.Coord, 0, len(probesMap))
	isAgain := false

	for coord, _ := range probesMap {
		probes = append(probes, coord)
	}

	route, err := c.client.GetRoutes(from, to, probes, c.setting.AvoidAreaOffset)
	if err != nil {
		return avoidAreas, avoidProbes, debug, err
	}
	routePoints := route[0].Points

	// Debug Info
	debugLog.Debug(count, routePoints, probes)

	debugCurToNextToProbe := make([]string, 0)
	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]
		probePoints := c.probeManager.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint probe.Probe) {
			if c.isAvoid(cur, next, probePoint) {
				curToNextToProbeInfo := fmt.Sprintf("needAvoid,fmt(cur;next;probeManager): %s", drive.FmtCoord(cur, next, probePoint.Coord))
				log.Info(curToNextToProbeInfo)
				debugCurToNextToProbe = append(debugCurToNextToProbe, curToNextToProbeInfo)
				probesMap[probePoint.Coord] = struct{}{}
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
		return avoidAreas, avoidProbes, debug, apis.AvoidOutOfRange
	}

	avoidAreas = make([][]drive.Coord, 0)
	avoidProbes = make([]drive.Coord, 0)
	for p := range probesMap {
		avoidAreas = append(avoidAreas, drive.ConvCoordToAvoidArea(p, c.setting.AvoidAreaOffset))
		avoidProbes = append(avoidProbes, p)
	}
	log.WithField("avoidPoints", probesMap).WithField("count", count).Info("根据直线距离计算需要避让的探头")
	return avoidAreas, avoidProbes, debug, nil
}

func (c *LBS) isAvoid(cur drive.Coord, next drive.Coord, probePoint probe.Probe) bool {

	//A1 -> A2   直线距离 B1
	//A1 -> 探头1 直线距离 B2
	//探头1 -> A2 直线距离 B3
	b1 := cur.GeoPoint().GreatCircleDistance(next.GeoPoint())
	b2 := cur.GeoPoint().GreatCircleDistance(probePoint.GeoPoint())
	b3 := probePoint.GeoPoint().GreatCircleDistance(next.GeoPoint())
	//B1 >= B2+B3-offset 即路过探头
	gap := b1 - (b2 + b3 - (float64(c.setting.Offset) / 1000))
	log.Infof("calculate: b1:%f, b2:%f, b3:%f, offset:%f, gap:%f", b1, b2, b3, float64(c.setting.Offset)/1000, gap)

	if gap <= 0 {
		return false
	}

	if len(probePoint.Towards) == 0 {
		return true
	}

	// A1->A2             朝向toward1
	// 探头->探头Towards    朝向toward2
	// toward1 与 toward2 差值小于 TowardRange 则视为同一个方向
	toward1 := cur.GeoPoint().BearingTo(next.GeoPoint())
	for _, toward := range probePoint.Towards {
		toward2 := probePoint.GeoPoint().BearingTo(toward.GeoPoint())
		towardGap := math.Abs(toward1 - toward2)
		if towardGap < c.setting.TowardRange || 360-towardGap < c.setting.TowardRange {
			return true
		}
	}
	return false
}

func (c *LBS) Probes(cur drive.Coord, near float64) ([]probe.Probe, error) {
	if near == 0 {
		return c.probeManager.Probes, nil
	}
	return c.probeManager.Near(cur, near), nil
}
