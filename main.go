package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/senthilsweb/ai-dlc-fp-estimation/server/handler"
)

//go:embed app
var appFS embed.FS

//go:embed data
var dataFS embed.FS

func main() {
	port := flag.String("port", getEnv("FP_PORT", "8080"), "Listen port")
	appName := flag.String("app", getEnv("FP_APP", "ai-agents-provly"), "Default dataset (data/<name>/) to serve at startup")
	logLevel := flag.String("log-level", getEnv("FP_LOG_LEVEL", "info"), "Log level (debug, info, warn, error)")
	logFormat := flag.String("log-format", getEnv("FP_LOG_FORMAT", "text"), "Log format (text, json)")
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %s", *logLevel)
	}
	log.SetLevel(level)
	if *logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	}

	appSub, err := fs.Sub(appFS, "app")
	if err != nil {
		log.Fatalf("Failed to load embedded app assets: %v", err)
	}
	dataSub, err := fs.Sub(dataFS, "data")
	if err != nil {
		log.Fatalf("Failed to load embedded data: %v", err)
	}

	if _, err := fs.Stat(dataSub, *appName); err != nil {
		log.Fatalf("Configured dataset %q not found under data/: %v", *appName, err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logrusMiddleware())

	r.GET("/api/data", handler.DataHandler(dataSub, *appName))
	r.GET("/api/apps", handler.AppsHandler(dataSub))
	r.NoRoute(handler.SPAHandler(appSub, ""))

	addr := ":" + *port
	log.WithFields(log.Fields{"addr": addr, "dataset": *appName}).Info("Server started")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func logrusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.WithFields(log.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
		}).Info("request")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
