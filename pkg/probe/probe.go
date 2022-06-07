package probe

import (
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/xyctruth/stream"
)

type ProbeManager struct {
	Probes []Probe `yaml:"probes"`
}

type Probe struct {
	drive.Coord `mapstructure:",squash"`
	Towards     []Toward `yaml:"towards" json:"towards"` // 探头朝向的坐标
}

type Toward struct {
	drive.Coord `mapstructure:",squash"`
	Value       float64 `json:"value" yaml:"value"` // 探头朝向的角度值
}

// Near 获取附近的探头
// distanceFilter 单位:KM
func (p *ProbeManager) Near(cur drive.Coord, distanceFilter float64) []Probe {
	return stream.NewSlice(p.Probes).Filter(func(point Probe) bool {
		distance := cur.GeoPoint().GreatCircleDistance(point.GeoPoint())
		return distance < distanceFilter
	}).ToSlice()
}

// All 全部探头
func (p *ProbeManager) All() []Probe {
	return p.Probes
}

func (p *ProbeManager) CalculateToward() {
	for _, probe := range p.Probes {
		for i, _ := range probe.Towards {
			probe.Towards[i].Value = probe.GeoPoint().BearingTo(probe.Towards[i].GeoPoint())
		}
	}
}

type ProbeSet map[drive.Coord]struct{}

func (s ProbeSet) ToSlice() []drive.Coord {
	probes := make([]drive.Coord, 0, len(s))
	for coord, _ := range s {
		probes = append(probes, coord)
	}
	return probes
}

func NewProbeSet() ProbeSet {
	return make(map[drive.Coord]struct{})
}
