package lbs

import (
	"fmt"
	"github.com/jftuga/geodist"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/xyctruth/stream"
)

// Calculator 根据路面距离计算需要避让的探头
type Calculator struct {
	client          drive.Client
	offset          int     // 路面距离计算偏移量，单位米
	avoidAreaOffset float64 // 生成四边形避让区偏移量, 单位经纬度
}

// NewCalculator
// client lbs client
// offset 路面距离计算偏移量，单位米
// avoidAreaOffset 生成四边形避让区偏移量, 单位经纬度
func NewCalculator(client drive.Client, offset int, avoidAreaOffset float64) *Calculator {
	c := &Calculator{
		client:          client,
		offset:          offset,
		avoidAreaOffset: avoidAreaOffset,
	}
	return c
}

// AvoidProbeByRoad 根据路面距离计算需要避让的探头
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
func (c *Calculator) AvoidProbeByRoad(from, to drive.Coord) (map[drive.Coord]struct{}, error) {
	avoidsMap := make(map[drive.Coord]struct{}, 0)
	count := 0

Again:
	count++
	avoids := make([]drive.Coord, 0, len(avoidsMap))
	isAgain := false

	for coord, _ := range avoidsMap {
		avoids = append(avoids, coord)
	}

	route, err := c.client.GetRoutes(from, to, avoids, c.avoidAreaOffset)
	if err != nil {
		return nil, err
	}
	routePoints := route[0].Points
	fmt.Println(drive.FmtCoord(routePoints...))
	// 需要避让的区域

	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]

		// 获取附近 5km 的探头
		probePoints := G_Probe.Near(cur, 5)
		if len(probePoints) == 0 {
			continue
		}

		curToNextAndProbes, err := c.client.GetDistanceMatrix([]drive.Coord{cur}, append([]drive.Coord{next}, probePoints...))
		if err != nil {
			return nil, err
		}

		probesToNext, err := c.client.GetDistanceMatrix(probePoints, []drive.Coord{next})
		if err != nil {
			return nil, err
		}

		// offset 单位米
		stream.NewSlice(probePoints).ForEach(func(i int, probePoint drive.Coord) {
			b1 := curToNextAndProbes[0]
			b2 := curToNextAndProbes[i+1]
			b3 := probesToNext[i]
			if b1 >= b2+b3-c.offset {
				fmt.Printf("needAvoid: %s  \n", drive.FmtCoord(cur, next, probePoint))
				fmt.Printf("needAvoid: b1:%d, b2:%d, b3:%d offset:%d  \n", b1, b2, b3, c.offset)
				avoidsMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}
	if isAgain {
		fmt.Println("again")
		goto Again
	}
	fmt.Println("执行次数：", count)
	return avoidsMap, nil
}

//AvoidProbeByLine 根据直线距离半径计算需要避让的探头
func (c *Calculator) AvoidProbeByLine(from, to drive.Coord) (map[drive.Coord]struct{}, error) {
	avoidsMap := make(map[drive.Coord]struct{}, 0)
	count := 0
Again:
	count++
	avoids := make([]drive.Coord, 0, len(avoidsMap))
	isAgain := false

	for coord, _ := range avoidsMap {
		avoids = append(avoids, coord)
	}

	route, err := c.client.GetRoutes(from, to, avoids, c.avoidAreaOffset)
	if err != nil {
		return nil, err
	}
	routePoints := route[0].Points

	for i := 0; i < len(routePoints)-1; i++ {
		cur := routePoints[i]
		next := routePoints[i+1]
		probePoints := G_Probe.All()
		stream.NewSlice(probePoints).ForEach(func(_ int, probePoint drive.Coord) {
			_, distance, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: next.Lat, Lon: next.Lon})
			_, probeDistance, _ := geodist.VincentyDistance(geodist.Coord{Lat: cur.Lat, Lon: cur.Lon}, geodist.Coord{Lat: probePoint.Lat, Lon: probePoint.Lon})
			if probeDistance < distance {
				fmt.Printf("needAvoid: %s  \n", drive.FmtCoord(cur, next, probePoint))
				avoidsMap[probePoint] = struct{}{}
				isAgain = true
			}
		})
	}
	if isAgain {
		fmt.Println("again")
		goto Again
	}
	fmt.Println("执行次数：", count)
	return avoidsMap, nil
}
