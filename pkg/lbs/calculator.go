package lbs

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/xyctruth/stream"
)

// Calculator 根据路面距离计算需要避让的探头
type Calculator struct {
	client drive.Client
}

func NewCalculator(client drive.Client) *Calculator {
	c := &Calculator{
		client: client,
	}
	return c
}

// AvoidProbeByRoad 根据直线距离半径计算需要避让的探头
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
func (c *Calculator) AvoidProbeByRoad(from, to drive.Coord, probePoints []drive.Coord) (map[drive.Coord]struct{}, error) {
	route, err := c.client.GetRoutes(from, to)
	if err != nil {
		panic(err)
	}
	routePoints := route[0].Points

	// 需要避让的区域
	avoidPoints := make(map[drive.Coord]struct{}, 0)

	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]

		// TODO 优化减少api请求数量
		curToNext, err := c.client.GetDistanceMatrix([]drive.Coord{cur}, []drive.Coord{next})
		if err != nil {
			return nil, err
		}

		curToProbes, err := c.client.GetDistanceMatrix([]drive.Coord{cur}, probePoints)
		if err != nil {
			return nil, err
		}

		probesToNext, err := c.client.GetDistanceMatrix(probePoints, []drive.Coord{next})
		if err != nil {
			return nil, err
		}

		// offset 单位米
		offset := 100
		stream.NewSlice(probePoints).ForEach(func(i int, probePoint drive.Coord) {
			b1 := curToNext[0]
			b2 := curToProbes[i]
			b3 := probesToNext[i]

			fmt.Printf("b1:%d , b2:%d, b3:%d  , isAvoid:%v  \n", b1, b2, b3, b1 >= b2+b3-offset)
			if b1 >= b2+b3-offset {
				avoidPoints[probePoint] = struct{}{}
			}
		})

	}
	return avoidPoints, nil
}

//AvoidProbeByLine 根据直线距离半径计算需要避让的探头
func (c *Calculator) AvoidProbeByLine(from, to drive.Coord, probePoints []drive.Coord) (map[drive.Coord]struct{}, error) {
	route, err := c.client.GetRoutes(from, to)
	if err != nil {
		panic(err)
	}
	routePoints := route[0].Points

	// 需要避让的区域
	avoidPoints := make(map[drive.Coord]struct{}, 0)

	for i := 0; i < len(routePoints)-1; i++ {
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint drive.Coord) {
			cur := routePoints[i]
			next := routePoints[i+1]
			_, distance, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
			_, probeDistance, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon})
			if probeDistance < distance {
				avoidPoints[probePoint] = struct{}{}
			}
		})
	}
	return avoidPoints, nil
}
