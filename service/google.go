package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	geojson "github.com/paulmach/go.geojson"
	"github.com/twpayne/go-polyline"
	"log"
	"math"
)

type Google struct {
	httpClient *resty.Client
}

type Response struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Legs             []Leg `json:"legs"`
	OverviewPolyline struct {
		Points string `json:"points"`
	} `json:"overview_polyline"`
}

type Leg struct {
	Distance struct {
		Text  string  `json:"text"`
		Value float32 `json:"value"`
	} `json:"distance"`
}

func NewGoogleService() Google {
	return Google{
		httpClient: resty.New(),
	}
}

func (s Google) Route(origin, destination [2]float64, mode string) (float32, geojson.FeatureCollection) {
	var res *resty.Response
	params := map[string]string{
		"origin":       fmt.Sprintf("%f,%f", origin[0], origin[1]),
		"destination":  fmt.Sprintf("%f,%f", destination[0], destination[1]),
		"mode":         mode,
		"key":          "AIzaSyAh8254LV3hM66lOdtuQwaLH1G8lDdPGDQ",
		"alternatives": "true",
	}
	res, err := s.httpClient.R().
		SetQueryParams(params).
		Get("https://maps.googleapis.com/maps/api/directions/json")
	if err != nil {
		return 0, geojson.FeatureCollection{}
	}
	var r Response
	if res.IsSuccess() {
		err = json.Unmarshal(res.Body(), &r)
		log.Println(err)
		if err != nil {
			return 0, geojson.FeatureCollection{}
		}
		route := shortestRoute(r.Routes)
		buf := []byte(route.OverviewPolyline.Points)
		coords, _, _ := polyline.DecodeCoords(buf)
		rightCoords := make([][]float64, 0)
		for _, c := range coords {
			rightCoords = append(rightCoords, []float64{c[1], c[0]})
		}
		geoj := geojson.NewFeatureCollection()
		geoj.AddFeature(geojson.NewLineStringFeature(rightCoords))
		return route.Legs[0].Distance.Value, *geoj
	}
	return 0, geojson.FeatureCollection{}
}

func shortestRoute(routes []Route) Route {
	lowestDistance := float32(math.MaxFloat32)
	result := Route{}
	for _, r := range routes {
		if r.Legs[0].Distance.Value < lowestDistance {
			result = r
			lowestDistance = r.Legs[0].Distance.Value
		}
	}
	return result
}
