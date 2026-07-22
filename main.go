package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"
	"strings"
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
	devMode := flag.Bool("dev", getEnv("FP_DEV", "false") == "true", "Serve app/ and data/ live from disk instead of the embedded copies — no rebuild needed to see edits (run from the repo root)")
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

	var appSub, dataSub fs.FS
	if *devMode {
		// Live filesystem, not the compiled-in copy: edit app/index.html or any
		// data/<app>/*.json and the next request picks it up — no `go build`.
		appSub = os.DirFS("app")
		dataSub = os.DirFS("data")
		if _, err := fs.Stat(appSub, "index.html"); err != nil {
			log.Fatalf("--dev requires an app/index.html relative to the current directory — run from the repo root: %v", err)
		}
	} else {
		appSub, err = fs.Sub(appFS, "app")
		if err != nil {
			log.Fatalf("Failed to load embedded app assets: %v", err)
		}
		dataSub, err = fs.Sub(dataFS, "data")
		if err != nil {
			log.Fatalf("Failed to load embedded data: %v", err)
		}
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
	log.WithFields(log.Fields{"addr": addr, "dataset": *appName, "dev": *devMode}).Info("Server started")
	if err := r.Run(addr); err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Fatalf("Port %s is already in use — pick another with --port <n> or FP_PORT=<n> (e.g. FP_PORT=8081): %v", *port, err)
		}
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
