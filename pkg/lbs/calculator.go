package lbs

import (
	"fmt"
	"github.com/jftuga/geodist"
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

//AvoidProbeByLine 根据直线距离计算需要避让的探头
func (c *Calculator) AvoidProbeByLine(from, to drive.Coord) (map[drive.Coord]struct{}, error) {
	logEntry := log.WithField("route", "AvoidProbeByLine")

	avoidsMap := make(map[drive.Coord]struct{}, 0)
	count := 0
Again:
	count++
	if count > c.setting.MaxRoute {
		return nil, fmt.Errorf("超出最大路线规划次数:%d", c.setting.MaxRoute)
	}

	avoids := make([]drive.Coord, 0, len(avoidsMap))
	isAgain := false

	for coord, _ := range avoidsMap {
		avoids = append(avoids, coord)
	}

	route, err := c.client.GetRoutes(from, to, avoids, c.setting.AvoidAreaOffset)
	if err != nil {
		return nil, err
	}
	routePoints := route[0].Points
	logEntry.Info("routes:", drive.FmtCoord(routePoints...))

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
			gap := b1 - (b2 + b3 - (float64(c.setting.LineOffset) / 1000))
			logEntry.Infof("calculate: b1:%f, b2:%f, b3:%f, offset:%f, gap:%f", b1, b2, b3, float64(c.setting.LineOffset)/1000, gap)
			if gap > 0 {
				logEntry.Infof("needAvoid: %s ", drive.FmtCoord(cur, next, probePoint))
				avoidsMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}
	if isAgain {
		logEntry.Info("again,count:", count)
		goto Again
	}
	if len(avoidsMap) > c.setting.MaxAvoid {
		return nil, fmt.Errorf("超出最大避让点：CurAvoid:%d MaxAvoid:%d", len(avoidsMap), c.setting.MaxAvoid)
	}

	logEntry.Info("执行次数：", count)
	return avoidsMap, nil
}

// AvoidProbeByRoad 根据路面距离计算需要避让的探头
//首先获取到路线A坐标串
//A1 -> A2 -> A3 -> A4 -> AN
//
//使用路面距离api判断是否经过探头
//api：https://lbs.qq.com/service/webService/webServiceGuide/webServiceMatrix
//
//A1 -> A2 路面距离 B1
//A1 -> 探头1 路面距离 B2
//探头1 -> A2 路面距离 B3
//
//B1 >= B2+B3 即路过探头
func (c *Calculator) AvoidProbeByRoad(from, to drive.Coord) (map[drive.Coord]struct{}, error) {
	logEntry := log.WithField("route", "AvoidProbeByRoad")
	avoidsMap := make(map[drive.Coord]struct{}, 0)
	count := 0

Again:
	count++
	if count > c.setting.MaxRoute {
		return nil, fmt.Errorf("超出最大路线规划次数:%d", c.setting.MaxRoute)
	}
	avoids := make([]drive.Coord, 0, len(avoidsMap))
	isAgain := false

	for coord, _ := range avoidsMap {
		avoids = append(avoids, coord)
	}

	route, err := c.client.GetRoutes(from, to, avoids, c.setting.AvoidAreaOffset)
	if err != nil {
		return nil, err
	}
	routePoints := route[0].Points
	logEntry.Info("routes:", drive.FmtCoord(routePoints...))
	// 需要避让的区域

	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]

		// 获取附近 5km 的探头
		probePoints := c.probe.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}

		curToNextAndProbes, err := c.client.GetDistanceMatrix([]drive.Coord{cur}, append([]drive.Coord{next}, probePoints...))
		if err != nil {
			return nil, err
		}

		probesToNext, err := c.client.GetDistanceMatrix(probePoints, []drive.Coord{next})
		if err != nil {
			return nil, err
		}

		// offset 单位米
		stream.NewSlice(probePoints).ForEach(func(i int, probePoint drive.Coord) {
			b1 := curToNextAndProbes[0]
			b2 := curToNextAndProbes[i+1]
			b3 := probesToNext[i]
			gap := b1 - (b2 + b3 - c.setting.RoadOffset)
			logEntry.Infof("calculate: b1:%d, b2:%d, b3:%d, offset:%d, gap:%d", b1, b2, b3, c.setting.RoadOffset, gap)
			if gap > 0 {
				logEntry.Infof("needAvoid: %s ", drive.FmtCoord(cur, next, probePoint))
				avoidsMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}
	if isAgain {
		logEntry.Info("again,count:", count)
		goto Again
	}
	logEntry.Info("执行次数：", count)
	return avoidsMap, nil
}
