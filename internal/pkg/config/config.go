package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               string
	StorageFile        string
	BaseURL            string
	ReadTimeout        int
	WriteTimeout       int
	Debug              bool
	LogLevel           string
	LogFormat          string
	RateLimit          int
	CorsAllowedOrigins []string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	storageFile := os.Getenv("STORAGE_FILE")
	if storageFile == "" {
		storageFile = "./storage.json"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	readTimeout := 30
	if val := os.Getenv("READ_TIMEOUT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			readTimeout = v
		}
	}

	writeTimeout := 30
	if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			writeTimeout = v
		}
	}

	debug := false
	if val := os.Getenv("DEBUG"); val != "" {
		debug = strings.ToLower(val) == "true"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "json"
	}

	rateLimit := 10
	if val := os.Getenv("RATE_LIMIT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			rateLimit = v
		}
	}

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var corsAllowedOrigins []string
	if corsOrigins == "" || corsOrigins == "*" {
		corsAllowedOrigins = []string{"*"}
	} else {
		corsAllowedOrigins = strings.Split(corsOrigins, ",")
	}

	return &Config{
		Port:               ":" + port,
		StorageFile:        storageFile,
		BaseURL:            baseURL,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		Debug:              debug,
		LogLevel:           logLevel,
		LogFormat:          logFormat,
		RateLimit:          rateLimit,
		CorsAllowedOrigins: corsAllowedOrigins,
	}
}
