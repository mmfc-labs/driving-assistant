package lbs

import (
	"fmt"
	geo "github.com/kellydunn/golang-geo"
	"testing"
)

func TestBearing(t *testing.T) {
	// 探头1 坐标
	a1 := geo.NewPoint(23.265333, 113.354559)
	// 探头1 朝向的方向 的坐标
	//a2 := geo.NewPoint(23.265076, 113.354948)

	a3 := geo.NewPoint(23.265136, 113.354565)
	fmt.Println(a1.BearingTo(a3))
	// output：125.72151935048862

}
