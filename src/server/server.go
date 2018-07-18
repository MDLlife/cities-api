package server

import (
	"github.com/boltdb/bolt"
	"cities-api/src/cache"
	"cities-api/src/config"
	"net/http"
	"time"
)

func Server(db *bolt.DB, options *config.Options, c *cache.Cache) *http.Server {
	return &http.Server{
		Addr:           ":" + options.Port,
		Handler:        NewRouter(db, options, c),
		ReadTimeout:    time.Duration(options.Timeout) * time.Second,
		WriteTimeout:   time.Duration(options.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
