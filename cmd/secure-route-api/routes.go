package main

import (
	graph "github.com/JesseleDuran/secure-route-api"
	handler "github.com/JesseleDuran/secure-route-api/http"
	"github.com/JesseleDuran/secure-route-api/s3"
	"github.com/JesseleDuran/secure-route-api/service"
	"github.com/gin-gonic/gin"
)

// setupRoutes returns a Gin server ready to rise up with all the available endpoints.
func setupRoutes(dataSource graph.Source, service service.Google) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(
		gin.Recovery(),
	)
	prefixv1 := router.Group("api/secure-route/v1")
	crimes := s3.NewCrimesInMemoryProvider(s3.GetClient())
	pathHandler := handler.PathHandler{GraphSource: dataSource, Google: service, Crimes: crimes.Fetch()}
	prefixv1.GET("path", pathHandler.ServeHTTP)

	return router
}
