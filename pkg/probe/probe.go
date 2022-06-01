package probe

import (
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/xyctruth/stream"
)

type Probe struct {
	Points []drive.Coord
}

// Near 获取附近的探头
// distanceFilter 单位:KM
func (p *Probe) Near(cur drive.Coord, distanceFilter float64) []drive.Coord {
	return stream.NewSlice(p.Points).Filter(func(point drive.Coord) bool {
		_, distance, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: point.Lat, Lon: point.Lon})
		return distance < distanceFilter
	}).ToSlice()
}

// All 全部探头
func (p *Probe) All() []drive.Coord {
	return p.Points
}
