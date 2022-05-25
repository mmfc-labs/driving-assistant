package tencent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
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

func (c *Client) GetRoutes(from, to drive.Coord) ([]drive.Route, error) {
	resp, err := c.httpClient.R().
		SetQueryParams(map[string]string{
			"from":     fmt.Sprintf("%f,%f", from.Lat, from.Lon),
			"to":       fmt.Sprintf("%f,%f", to.Lat, to.Lon),
			"output":   "json",
			"callback": "cb",
			"key":      c.key,
		}).
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
		return nil, errors.New("Routes:" + p.Message)
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
	fmt.Println(p)

	if p.Status > 0 {
		return nil, errors.New("DistanceMatrix:" + p.Message)
	}

	result := make([]int, 0, 0)
	for _, row := range p.Result.Rows {
		for _, element := range row.Elements {
			result = append(result, element.Distance)
		}
	}

	// TODO 暂时解决每秒请求量已达到上限
	time.Sleep(time.Millisecond * 1000)
	return result, err
}