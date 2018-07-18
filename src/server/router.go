package server

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"cities-api/src/cache"
	"cities-api/src/config"
	"cities-api/src/middleware"
)

func NewRouter(db *bolt.DB, options *config.Options, c *cache.Cache) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CommonHeaders(options.CORSOrigins))

	v1 := router.Group("/1.0")
	{
		v1.GET("/application/status", MakeApplicationStatusEndpoint(db))
		v1.GET("/cities/:id", MakeCityEndpoint(db, options))
		v1.GET("/cities", MakeCitiesEndpoint(db, options, c))
	}

	return router
}
