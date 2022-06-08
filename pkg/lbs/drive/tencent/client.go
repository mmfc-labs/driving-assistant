package tencent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/stream"
	"strings"
	"time"
)

type Client struct {
	key        string
	httpClient *resty.Client
}

func NewClient(key string) *Client {
	httpClient := resty.New()
	httpClient.SetBaseURL("https://apis.map.qq.com/ws")
	c := &Client{
		key:        key,
		httpClient: httpClient,
	}
	return c
}

func (c *Client) GetRoutes(from, to drive.Coord, avoids []drive.Coord, avoidAreaOffset float64) ([]drive.Route, error) {
	reduceFunc := func(r string, e drive.Coord) string {
		avoidAreas := drive.ConvCoordToAvoidArea(e, avoidAreaOffset)
		r += strings.Trim(stream.NewSliceByMapping[drive.Coord, string, string](avoidAreas).Reduce(func(r string, e drive.Coord) string {
			r += fmt.Sprintf("%f,%f;", e.Lat, e.Lon)
			return r
		}), ";") + "|"
		return r
	}
	avoidPolygonsParam := strings.Trim(stream.NewSliceByMapping[drive.Coord, string, string](avoids).Reduce(reduceFunc), "|")
	params := map[string]string{
		"from":     fmt.Sprintf("%f,%f", from.Lat, from.Lon),
		"to":       fmt.Sprintf("%f,%f", to.Lat, to.Lon),
		"output":   "json",
		"callback": "cb",
		"key":      c.key,
	}
	if avoidPolygonsParam != "" {
		params["avoid_polygons"] = avoidPolygonsParam
	}
Retry:
	resp, err := c.httpClient.R().
		SetQueryParams(params).
		Get("/direction/v1/driving/")

	if err != nil {
		return nil, err
	}

	p := struct {
		Status  int
		Message string
		Result  struct {
			Routes []struct {
				Polyline []float64
			}
		}
	}{}

	err = json.Unmarshal(resp.Body(), &p)
	if err != nil {
		return nil, err
	}

	if p.Status > 0 {
		log.Errorf("tencent route status:%d , message:%s", p.Status, p.Message)
		if p.Status == 120 {
			time.Sleep(time.Millisecond * 500)
			goto Retry
		}
		if p.Status == 410 {
			return nil, apis.ErrorRouteFailed
		}
		return nil, errors.New("RouteLogs:" + p.Message)
	}
	routes := make([]drive.Route, 0, len(p.Result.Routes))

	for _, route := range p.Result.Routes {
		for i := 2; i < len(route.Polyline); i++ {
			route.Polyline[i] = route.Polyline[i-2] + route.Polyline[i]/1000000
		}
		points := make([]drive.Coord, 0, len(route.Polyline)/2)
		for i := 0; i < len(route.Polyline); i = i + 2 {
			points = append(points, drive.Coord{
				Lat: route.Polyline[i],
				Lon: route.Polyline[i+1],
			})

		}
		routes = append(routes, drive.Route{Points: points})

	}
	return routes, nil
}

func (c *Client) GetDistanceMatrix(froms, tos []drive.Coord) ([]int, error) {
	reduceFunc := func(r string, e drive.Coord) string {
		r += fmt.Sprintf("%f,%f;", e.Lat, e.Lon)
		return r
	}
	fromParam := strings.Trim(stream.NewSliceByMapping[drive.Coord, string, string](froms).Reduce(reduceFunc), ";")
	toParam := strings.Trim(stream.NewSliceByMapping[drive.Coord, string, string](tos).Reduce(reduceFunc), ";")

Retry:
	resp, err := c.httpClient.R().
		SetQueryParams(map[string]string{
			"from":     fromParam,
			"to":       toParam,
			"output":   "json",
			"callback": "cb",
			"key":      c.key,
			"mode":     "driving",
		}).
		Get("/distance/v1/matrix")

	if err != nil {
		return nil, err
	}

	p := struct {
		Status  int
		Message string
		Result  struct {
			Rows []struct {
				Elements []struct {
					Distance int //米
					Duration int //秒
				}
			}
		}
	}{}

	err = json.Unmarshal(resp.Body(), &p)
	if err != nil {
		return nil, err
	}

	if p.Status > 0 {
		if p.Status == 120 {
			time.Sleep(time.Millisecond * 500)
			goto Retry
		}

		return nil, errors.New("DistanceMatrix:" + p.Message)
	}

	result := make([]int, 0, 0)
	for _, row := range p.Result.Rows {
		for _, element := range row.Elements {
			result = append(result, element.Distance)
		}
	}
	return result, err
}
