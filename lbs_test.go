package main

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive/tencent"
	"testing"
)

func TestAvoidByLine(t *testing.T) {
	// 起点，终点
	from, to := drive.Coord{Lat: 22.577781, Lon: 113.910683}, drive.Coord{Lat: 22.576752, Lon: 113.914866}
	calculator := lbs.NewCalculator(tencent.NewClient(config.TencentKey), config.Offset, config.AvoidAreaOffset)

	//根据直线距离计算需要避让的探头
	avoidPoints, err := calculator.AvoidProbeByLine(from, to)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("根据直线距离计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}
}

func TestAvoidByRoad(t *testing.T) {
	// 起点，终点
	from, to := drive.Coord{Lat: 22.560447, Lon: 113.874653}, drive.Coord{Lat: 22.55453, Lon: 113.887378}
	calculator := lbs.NewCalculator(tencent.NewClient(config.TencentKey), config.Offset, config.AvoidAreaOffset)

	//根据路面距离计算需要避让的探头
	avoidPoints, err := calculator.AvoidProbeByRoad(from, to)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("根据路面距离计算需要避让的探头")
	for key, _ := range avoidPoints {
		fmt.Println(key)
	}
}

func TestMaxPolyline(t *testing.T) {
	client := tencent.NewClient(config.TencentKey)
	from, to := drive.Coord{Lat: 22.575098, Lon: 113.85605}, drive.Coord{Lat: 22.55453, Lon: 113.887378}
	route, err := client.GetRoutes(from, to, nil, 0)
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
