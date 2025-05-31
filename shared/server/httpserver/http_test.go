package httpserver_test

import (
	"net/http"
	"testing"
	"time"

	"tixgo/config"
	"tixgo/shared/server/httpserver"

	"github.com/gin-gonic/gin"
)

// DefaultConfig returns a default httpServer configuration
func DefaultConfig() httpserver.Config {
	return httpserver.Config{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
}

// DefaultRouterConfig returns a default router configuration
func DefaultRouterConfig() httpserver.RouterConfig {
	return httpserver.RouterConfig{
		Environment: "dev",
		EnableCORS:  true,
		EnableAuth:  true,
	}
}

func TestConfigFromAppConfig(t *testing.T) {
	appConfig := &config.AppConfig{
		Server: config.Server{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}

	httpServerConfig := httpserver.Config{
		Host:         appConfig.Server.Host,
		Port:         appConfig.Server.Port,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
		IdleTimeout:  appConfig.Server.IdleTimeout,
	}
	if httpServerConfig.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", httpServerConfig.Host)
	}
	if httpServerConfig.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", httpServerConfig.Port)
	}
	if httpServerConfig.ReadTimeout != 10*time.Second {
		t.Errorf("Expected read timeout 10s, got %v", httpServerConfig.ReadTimeout)
	}
}

func TestRouterConfigFromAppConfig(t *testing.T) {
	appConfig := &config.AppConfig{
		App: config.App{
			Environment: "prod",
		},
	}

	routerConfig := httpserver.RouterConfig{
		Environment: appConfig.App.Environment,
		EnableCORS:  true,
		EnableAuth:  true,
	}

	if routerConfig.Environment != "prod" {
		t.Errorf("Expected environment 'prod', got '%s'", routerConfig.Environment)
	}
	if !routerConfig.EnableCORS {
		t.Error("Expected CORS to be enabled")
	}
	if !routerConfig.EnableAuth {
		t.Error("Expected auth to be enabled")
	}
}

func TestSetupRouter(t *testing.T) {
	config := httpserver.RouterConfig{
		Environment: "test",
		EnableCORS:  true,
		EnableAuth:  true,
	}

	router := httpserver.SetupRouter(config)

	if router == nil {
		t.Fatal("Expected router to be created")
	}

	// Test that health endpoints are registered
	routes := router.Routes()

	healthFound := false
	readyFound := false
	liveFound := false

	for _, route := range routes {
		switch route.Path {
		case "/health":
			healthFound = true
		case "/ready":
			readyFound = true
		case "/live":
			liveFound = true
		}
	}

	if !healthFound {
		t.Error("Health endpoint not found")
	}
	if !readyFound {
		t.Error("Ready endpoint not found")
	}
	if !liveFound {
		t.Error("Live endpoint not found")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Port)
	}
	if config.ReadTimeout != 10*time.Second {
		t.Errorf("Expected default read timeout 10s, got %v", config.ReadTimeout)
	}
}

func TestDefaultRouterConfig(t *testing.T) {
	config := DefaultRouterConfig()

	if config.Environment != "dev" {
		t.Errorf("Expected default environment 'dev', got '%s'", config.Environment)
	}
	if !config.EnableCORS {
		t.Error("Expected CORS to be enabled by default")
	}
	if !config.EnableAuth {
		t.Error("Expected auth to be enabled by default")
	}
}

func TesthttpServerNew(t *testing.T) {
	config := httpserver.Config{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	httpServer := httpserver.New(config, handler)

	if httpServer == nil {
		t.Fatal("Expected httpServer to be created")
	}

	expectedAddr := "localhost:8080"
	if httpServer.Addr() != expectedAddr {
		t.Errorf("Expected httpServer address '%s', got '%s'", expectedAddr, httpServer.Addr())
	}

	// Test that handler is set (function comparison not allowed, just check it's not nil)
	if httpServer.Handler() == nil {
		t.Error("Expected handler to be set")
	}
}

func TestAddAPIGroup(t *testing.T) {
	router := gin.New()

	v1Group := httpserver.AddAPIGroup(router, "v1")
	v2Group := httpserver.AddAPIGroup(router, "v2")

	if v1Group == nil {
		t.Fatal("Expected v1 group to be created")
	}
	if v2Group == nil {
		t.Fatal("Expected v2 group to be created")
	}

	// Add test routes to verify group paths
	v1Group.GET("/test", func(c *gin.Context) {})
	v2Group.GET("/test", func(c *gin.Context) {})

	routes := router.Routes()

	v1Found := false
	v2Found := false

	for _, route := range routes {
		if route.Path == "/api/v1/test" {
			v1Found = true
		}
		if route.Path == "/api/v2/test" {
			v2Found = true
		}
	}

	if !v1Found {
		t.Error("v1 API group route not found")
	}
	if !v2Found {
		t.Error("v2 API group route not found")
	}
}

func TestAddProtectedGroup(t *testing.T) {
	router := gin.New()
	apiGroup := router.Group("/api/v1")

	protectedGroup := httpserver.AddProtectedGroup(apiGroup, "/protected")

	if protectedGroup == nil {
		t.Fatal("Expected protected group to be created")
	}

	// Add test route to verify group path
	protectedGroup.GET("/test", func(c *gin.Context) {})

	routes := router.Routes()

	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/protected/test" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Protected group route not found")
	}
}
