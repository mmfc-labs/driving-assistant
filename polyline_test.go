package main

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/lbs"
	"github.com/mmfc-labs/driving-assistant/lbs/tencent"
	"testing"
)

const TencentKey = "KN6BZ-G526D-JAI4V-PGSJ2-6L5U6-YYFBV"

func TestGetAvoid(t *testing.T) {
	allAvoidPoints := []lbs.Coord{
		{
			Lat: 22.565615,
			Lon: 113.86821,
		},
		{
			Lat: 22.557223,
			Lon: 113.877889,
		},
		{
			Lat: 23.565615,
			Lon: 114.86821,
		},
	}
	needAvoidPoints := make(map[lbs.Coord]struct{}, 0)

	client := tencent.NewLBSClient(TencentKey)
	route, err := client.GetRoute(22.575098, 113.85605, 22.55453, 113.887378)
	if err != nil {
		panic(err)
	}
	points := route[0].Points
	for i := 0; i < len(points); i++ {
		for _, avoidPoint := range allAvoidPoints {
			_, km, _ := geodist.VincentyDistance(
				geodist.Coord{Lat: points[i].Lat, Lon: points[i].Lon},
				geodist.Coord{Lat: avoidPoint.Lat, Lon: avoidPoint.Lon})

			if km < 1 {
				needAvoidPoints[avoidPoint] = struct{}{}
			}
		}
	}

	for key, _ := range needAvoidPoints {
		fmt.Println("需要规避的点", key)
	}
}

func TestMaxPolyline(t *testing.T) {
	lbs := tencent.NewLBSClient(TencentKey)
	route, err := lbs.GetRoute(22.529293, 113.971425, 22.544021, 113.989136)
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
