package http

import (
	"context"
	"errors"
	api "github.com/JesseleDuran/secure-route-api"
	"github.com/JesseleDuran/secure-route-api/service"
	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrWrongOriginValues      = errors.New("invalid origin lat, lng")
	ErrWrongDestinationValues = errors.New("invalid destination lat, lng")
)

type PathHandler struct {
	GraphSource api.Source
	Google      service.Google
	Crimes      map[uint64]api.Crime
}

func (handler PathHandler) ServeHTTP(ctx *gin.Context) {
	origin := strings.Split(ctx.Query("origin"), ",")
	destination := strings.Split(ctx.Query("destination"), ",")
	modeS := ctx.Query("mode")
	content := ctx.Query("content")
	source := ctx.Query("source")

	if len(origin) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrWrongOriginValues.Error()})
		return
	}
	latOri, err := strconv.ParseFloat(origin[0], 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lngOri, err := strconv.ParseFloat(origin[1], 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(destination) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrWrongDestinationValues.Error()})
		return
	}
	latDest, err := strconv.ParseFloat(destination[0], 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lngDest, err := strconv.ParseFloat(destination[1], 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var crimesIDs []uint64
	d, path := float32(0), geojson.FeatureCollection{}
	if source == "google" {
		d, path = handler.Google.Route([2]float64{latOri, lngOri}, [2]float64{latDest, lngDest}, modeS)
	} else {
		mode, err := api.ModeFromString(modeS)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		repo, _ := handler.GraphSource.GraphByModeAndContent(mode, content)
		d, path, crimesIDs = repo.Path(context.Background(), [2]float64{latOri, lngOri}, [2]float64{latDest, lngDest}, repo.Index())
	}
	resCrimes := make([]api.Crime, 0)
	for i, id := range crimesIDs {
		resCrimes = append(resCrimes, handler.Crimes[id])
		if i == 15 {
			break
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"distance":     d,
		"total_crimes": len(crimesIDs),
		"path":         path,
		"crimes":       resCrimes,
	})
}
