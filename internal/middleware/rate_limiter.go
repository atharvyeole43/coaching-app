package middleware

import (
	"time"

	"github.com/Reugito/dynamicratelimiter/config"
	"github.com/Reugito/dynamicratelimiter/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRateLimiter(router *gin.Engine) {

	rateLimitConf := config.RateLimitConfig{
		Redis: config.RedisConfig{
			EnableRedis: false,
		},
		RateLimits: config.RateLimitSettings{
			GlobalMaxRequestsPerSec: 30,
			DefaultRequestsPerSec:   10,
			MonitoringTimeFrame:     1 * time.Minute,
			IncreaseFactor:          2,
			IPExceedThreshold:       2,
		},
		EnableAdaptiveRateLimit: true,
	}

	rl := middleware.NewRateLimiter(rateLimitConf)
	router.Use(rl.Middleware())

	router.GET("/rate-limit-conf", rl.RateLimitMetricsHandler())
	router.GET("/reset-limit-conf", rl.DefaultRequestsPerSec())

}
