package lbs

import (
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
)

type Avoid interface {
	Calculate(routePoints []drive.Coord, probePoints []drive.Coord) map[drive.Coord]struct{}
}

// AvoidByRoad 根据路面距离计算避让的区域
type AvoidByRoad struct {
}

func (avoid AvoidByRoad) Calculate(routePoints []drive.Coord, probePoints []drive.Coord) map[drive.Coord]struct{} {
	// 需要避让的区域
	avoidPoints := make(map[drive.Coord]struct{}, 0)
	return avoidPoints
}

// AvoidByStraightLine 根据直线距离半径计算避让的区域
type AvoidByStraightLine struct {
}

func (avoid AvoidByStraightLine) Calculate(routePoints []drive.Coord, probePoints []drive.Coord) map[drive.Coord]struct{} {
	// 需要避让的区域
	avoidPoints := make(map[drive.Coord]struct{}, 0)

	for i := 0; i < len(routePoints)-1; i++ {
		for _, avoidPoint := range probePoints {
			_, distance, _ := geodist.VincentyDistance(geodist.Coord{Lat: routePoints[i].Lat, Lon: routePoints[i].Lon}, geodist.Coord{Lat: routePoints[i+1].Lat, Lon: routePoints[i+1].Lon})
			_, km, _ := geodist.VincentyDistance(geodist.Coord{Lat: routePoints[i].Lat, Lon: routePoints[i].Lon}, geodist.Coord{Lat: avoidPoint.Lat, Lon: avoidPoint.Lon})
			if km < distance {
				avoidPoints[avoidPoint] = struct{}{}
			}
		}
	}
	return avoidPoints
}
