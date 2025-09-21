package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cctw-zed/wonder/internal/container"
	"github.com/cctw-zed/wonder/internal/server"
)

func main() {
	// Parse command line flags
	var configPath = flag.String("config", "", "Path to configuration file")
	var environment = flag.String("env", "", "Environment (development, testing, production)")
	flag.Parse()

	// Create application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize container with configuration
	var c *container.Container
	var err error

	if *environment != "" {
		// Load environment-specific configuration from configs directory
		c, err = container.NewContainerForEnvironment(ctx, *environment)
	} else if *configPath != "" {
		// Load configuration from custom path
		c, err = container.NewContainerWithConfig(ctx, *configPath)
	} else {
		// Load default configuration from configs directory
		c, err = container.NewContainer()
	}

	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Error closing container: %v", err)
		}
	}()

	// Create and start server
	srv := server.New(c)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s server on %s (environment: %s)",
			c.Config.App.Name,
			srv.GetAddr(),
			c.Config.App.Environment)

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
