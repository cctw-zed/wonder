package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cctw-zed/wonder/internal/container"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	container  *container.Container
}

// New creates a new server instance
func New(c *container.Container) *Server {
	// Set Gin mode based on environment
	switch c.Config.App.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Setup HTTP router
	router := setupRouter(c)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Config.Server.Host, c.Config.Server.Port),
		Handler:      router,
		ReadTimeout:  c.Config.Server.ReadTimeout,
		WriteTimeout: c.Config.Server.WriteTimeout,
		IdleTimeout:  c.Config.Server.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		container:  c,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	return s.httpServer.Addr
}

// setupRouter configures the HTTP routes
func setupRouter(c *container.Container) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS middleware if enabled
	if c.Config.Server.EnableCORS {
		router.Use(corsMiddleware())
	}

	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		// Check database health
		if err := c.Database.Health(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"app":         c.Config.App.Name,
			"version":     c.Config.App.Version,
			"environment": c.Config.App.Environment,
		})
	})

	// API version 1
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("/register", c.UserHandler.Register)
			// Add more user routes as needed
		}
	}

	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}