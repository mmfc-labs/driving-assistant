package lbs

import (
	"encoding/json"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/xyctruth/stream"
	"io/ioutil"
)

var (
	G_Probe *Probe
)

func init() {
	G_Probe = &Probe{}

	bytes, err := ioutil.ReadFile("./probe.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, G_Probe)
	if err != nil {
		panic(err)
	}
}

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
