package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/containous/traefik/log"

	"github.com/mashenjun/mirage/config"
	"github.com/mashenjun/mirage/pkg/endpoint"
	"github.com/mashenjun/mirage/pkg/middleware/prom"
	"github.com/mashenjun/mirage/pkg/service"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "cfgPath", "./", "intelliecard -cfgPath=/path/to/config.yaml")

	flag.Parse()
	if err := config.InitOption(cfgPath); err != nil {
		log.Panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	// router
	router := gin.New()

	// middleware
	router.Use(gin.Recovery())
	router.Use(prom.MetricsMiddleware)
	// router.Use(gzip.Gzip(gzip.BestCompression))
	pprof.Register(router)
	prom.Register(router)

	lvl, err := config.Options.Log.GetLevel()
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(lvl)
	fd, err := config.Options.Log.GetWriter()
	if err != nil {
		log.Panic(err)
	}
	defer fd.Close()
	log.SetOutput(fd)

	rdsCli := redis.NewClient(config.Options.RedisConfig)

	srv, err := service.New(rdsCli)
	if err != nil {
		log.Panic(err)
	}

	ep, err := endpoint.New(srv)
	if err != nil {
		log.Panic(err)
	}

	ep.MountOn(router)

	httpServer := &http.Server{
		Addr:              config.Options.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(config.Options.Server.RTimeout) * time.Second,
		ReadTimeout:       time.Duration(config.Options.Server.RTimeout) * time.Second,
		WriteTimeout:      time.Duration(config.Options.Server.WTimeout) * time.Second,
		IdleTimeout:       time.Duration(config.Options.Server.DTimeout) * time.Second,
		MaxHeaderBytes:    config.Options.Server.MaxHeaderBytes,
	}

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen err: %s", err)
		}
	}()
	log.Infof("listen on %s", config.Options.Server.Addr)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Infof("server shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown err:%+v", err)
	}
	log.Infof("server exiting")
}
