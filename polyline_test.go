package main

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"testing"
)

const TencentKey = "RB7BZ-CGUKW-AU4RO-RIFUW-57GFS-L4BP7"

func TestAvoidByRoad(t *testing.T) {
	// 探头
	calculator := lbs.NewCalculator(tencent.NewClient(TencentKey))
	from, to := drive.Coord{Lat: 22.560447, Lon: 113.874653}, drive.Coord{Lat: 22.55453, Lon: 113.887378}
	probe := lbs.LoadProbe()
	avoidPoints, err := calculator.AvoidProbeByRoad(from, to, probe.Points)
	if err != nil {
		panic(err)
	}
	fmt.Println("根据路面距离计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}

	// 根据直线距离半径计算需要避让的探头
	avoidPoints, err = calculator.AvoidProbeByLine(from, to, probe.Points)
	if err != nil {
		panic(err)
	}
	fmt.Println("根据直线距离半径计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}
}

func TestMaxPolyline(t *testing.T) {
	client := tencent.NewClient(TencentKey)
	from, to := drive.Coord{Lat: 22.575098, Lon: 113.85605}, drive.Coord{Lat: 22.55453, Lon: 113.887378}
	route, err := client.GetRoutes(from, to)
	if err != nil {
		panic(err)
	}

	points := route[0].Points
	var maxDistance float64
	for i := 0; i < len(points); i++ {
		fmt.Printf("%f,%f ", points[i].Lat, points[i].Lon)
		if i > 0 {
			var newYork = geodist.Coord{Lat: points[i-1].Lat, Lon: points[i-1].Lon}
			var sanDiego = geodist.Coord{Lat: points[i].Lat, Lon: points[i].Lon}
			_, km, _ := geodist.VincentyDistance(newYork, sanDiego)
			if maxDistance < km {
				maxDistance = km
			}
			fmt.Printf("距离上一个坐标 %f km\n", km)

		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Println("最长距离为：", maxDistance)
}
