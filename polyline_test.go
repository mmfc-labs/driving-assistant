package main

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"testing"
)

const TencentKey = "KN6BZ-G526D-JAI4V-PGSJ2-6L5U6-YYFBV"

func TestAvoidByRoad(t *testing.T) {
	// 探头
	probe := lbs.LoadProbe()

	client := tencent.NewClient(TencentKey)
	route, err := client.GetRoutes(drive.Coord{Lat: 22.575098, Lon: 113.85605}, drive.Coord{Lat: 22.55453, Lon: 113.887378})
	if err != nil {
		panic(err)
	}
	// 需要避让的区域
	avoid := lbs.AvoidByRoad{}
	avoidPoints := avoid.Calculate(route[0].Points, probe.Points)

	for key, _ := range avoidPoints {
		fmt.Println("需要规避的点", key)
	}
}

func TestAvoidByStraightLine(t *testing.T) {
	// 探头
	probe := lbs.LoadProbe()

	client := tencent.NewClient(TencentKey)
	route, err := client.GetRoutes(drive.Coord{Lat: 22.575098, Lon: 113.85605}, drive.Coord{Lat: 22.55453, Lon: 113.887378})
	if err != nil {
		panic(err)
	}
	// 需要避让的区域
	avoid := lbs.AvoidByStraightLine{}
	avoidPoints := avoid.Calculate(route[0].Points, probe.Points)

	for key, _ := range avoidPoints {
		fmt.Println("需要规避的点", key)
	}
}

func TestMaxPolyline(t *testing.T) {
	client := tencent.NewClient(TencentKey)
	route, err := client.GetRoutes(drive.Coord{Lat: 22.575098, Lon: 113.85605}, drive.Coord{Lat: 22.55453, Lon: 113.887378})
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
