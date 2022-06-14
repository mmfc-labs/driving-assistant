package lbs

import (
	"fmt"
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
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
	probeManager probe.Manager
}

// NewLBS
func NewLBS(setting config.Setting, probeManager probe.Manager) *LBS {
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
func (c *LBS) Route(from, to geo.Coord) (avoidAreas [][]geo.Coord, avoidProbes []geo.Coord, debug *apis.Debug, err error) {
	probesSet := probe.NewProbeSet()
	debug = &apis.Debug{}
	debug.RouteLogs = make([]*apis.DebugLog, 0)
	routeCount := 0

Again:
	routeCount++
	isAgain := false

	if routeCount > c.setting.MaxRoute {
		return avoidAreas, avoidProbes, debug, apis.ErrorRouteOutOfRange
	}

	route, err := c.client.GetRoutes(from, to, probesSet.ToSlice(), c.setting.AvoidAreaOffset)
	if err != nil {
		return avoidAreas, avoidProbes, debug, err
	}
	routePoints := route[0].Points

	// Add Debug Info
	debugLog := &apis.DebugLog{}
	debugCurToNextToProbe := make([]string, 0)
	debugLog.Debug(routeCount, routePoints, probesSet.ToSlice())

	// 循环坐标串
	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]
		//获取附近5km的探头
		probePoints := c.probeManager.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint probe.Probe) {
			if c.isAvoid(cur, next, probePoint) {
				curToNextToProbeInfo := fmt.Sprintf("第%d次路线规划,新增避让探头,格式为(当前坐标串A1;下一个坐标串A2;探头坐标): %s", routeCount, geo.Coords{cur, next, probePoint.Coord})
				log.Info(curToNextToProbeInfo)
				debugCurToNextToProbe = append(debugCurToNextToProbe, curToNextToProbeInfo)
				probesSet.Add(probePoint.Coord)
				isAgain = true
			}
		})
	}

	// Add Debug Info
	debug.Debug(debugLog, debugCurToNextToProbe, routeCount, probesSet.ToSlice())

	if probesSet.Len() > c.setting.MaxAvoid {
		return avoidAreas, avoidProbes, debug, apis.ErrorAvoidOutOfRange
	}

	if isAgain {
		log.Info("again,routeCount:", routeCount)
		goto Again
	}

	avoidAreas = probesSet.ToAvoidArea(c.setting.AvoidAreaOffset)
	avoidProbes = probesSet.ToSlice()
	log.Infof("路线规划完成,路线规划次数:%d,避让的探头:%s", routeCount, geo.Coords(avoidProbes))
	return avoidAreas, avoidProbes, debug, nil
}

// isAvoid 根绝 A1 A2 探头 来判断是否需要避让
func (c *LBS) isAvoid(cur geo.Coord, next geo.Coord, probePoint probe.Probe) bool {
	offsetKM := float64(c.setting.Offset) / 1000

	//A1 -> A2   直线距离 B1
	//A1 -> 探头1 直线距离 B2
	//探头1 -> A2 直线距离 B3
	b1 := cur.Distance(next)
	b2 := cur.Distance(probePoint.Coord)
	b3 := probePoint.Distance(next)
	//B1 >= B2+B3-offset 即路过探头
	gap := b1 - (b2 + b3)
	fmt.Println(gap)

	if gap+offsetKM <= 0 {
		return false
	}

	//判断三角形的高度
	s := (b1 + b2 + b3) / 2
	area := math.Sqrt(s * (s - b1) * (s - b2) * (s - b3))
	height := area * 2 / b1
	if height > offsetKM {
		return false
	}

	// 判断朝向
	if len(probePoint.Towards) == 0 {
		return true
	}

	// A1->A2             朝向toward1
	// 探头->探头Towards    朝向toward2
	// toward1 与 toward2 差值小于 TowardRange 则视为同一个方向
	toward1 := cur.BearingTo(next)
	for _, toward := range probePoint.Towards {
		toward2 := probePoint.BearingTo(toward.Coord)
		towardGap := math.Abs(toward1 - toward2)
		if towardGap < c.setting.TowardRange || 360-towardGap < c.setting.TowardRange {
			return true
		}
	}
	return false
}

func (c *LBS) Probes(cur geo.Coord, near float64) ([]probe.Probe, error) {
	if near == 0 {
		return c.probeManager.Probes, nil
	}
	return c.probeManager.Near(cur, near), nil
}
