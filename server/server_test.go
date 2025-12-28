package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test Server Creation and Configuration =============================================================================

// TestNew tests that New() creates a server with correct defaults
func TestNew(t *testing.T) {
	srv := New()

	if srv.port != 8080 {
		t.Errorf("Expected default port 8080, got %d", srv.port)
	}

	if srv.host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", srv.host)
	}

	if srv.routes == nil {
		t.Error("Expected routes to be initialized")
	}

	if len(srv.routes) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(srv.routes))
	}
}

// TestFluentAPI tests the fluent API pattern
func TestFluentAPI(t *testing.T) {
	srv := New().
		Port(3000).
		Host("localhost")

	if srv.port != 3000 {
		t.Errorf("Expected port 3000, got %d", srv.port)
	}

	if srv.host != "localhost" {
		t.Errorf("Expected host localhost, got %s", srv.host)
	}
}

// TestRoute tests adding routes
func TestRoute(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	srv := New().
		Route("/hello", handler).
		Route("/world", handler)

	if len(srv.routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(srv.routes))
	}

	if srv.routes[0].Path != "/hello" {
		t.Errorf("Expected first route path /hello, got %s", srv.routes[0].Path)
	}

	if srv.routes[1].Path != "/world" {
		t.Errorf("Expected second route path /world, got %s", srv.routes[1].Path)
	}
}

// Test Callbacks ==================================================================================================

// TestOnRequestCallback tests the OnRequest callback
func TestOnRequestCallback(t *testing.T) {
	called := false
	var capturedInfo RequestInfo

	handler := func(w http.ResponseWriter, r *http.Request) error {
		return Text(w, 200, "OK")
	}

	srv := New().
		Port(9876).
		Route("/test", handler).
		OnRequest(func(info RequestInfo) {
			called = true
			capturedInfo = info
		}).
		Build()

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call the handler directly (wrapped)
	srv.mux.ServeHTTP(w, req)

	if !called {
		t.Error("OnRequest callback was not called")
	}

	if capturedInfo.Method != "GET" {
		t.Errorf("Expected method GET, got %s", capturedInfo.Method)
	}

	if capturedInfo.Path != "/test" {
		t.Errorf("Expected path /test, got %s", capturedInfo.Path)
	}
}

// TestOnStartCallback tests the OnStart callback
func TestOnStartCallback(t *testing.T) {
	called := false

	srv := New().
		Port(9877).
		OnStart(func() {
			called = true
		}).
		Build()

	srv.Start()
	defer srv.Stop()

	// Give it a moment to call the callback
	time.Sleep(10 * time.Millisecond)

	if !called {
		t.Error("OnStart callback was not called")
	}
}

// TestOnStopCallback tests the OnStop callback
func TestOnStopCallback(t *testing.T) {
	called := false

	srv := New().
		Port(9878).
		OnStop(func() {
			called = true
		}).
		Build()

	srv.Start()
	time.Sleep(10 * time.Millisecond)
	srv.Stop()

	if !called {
		t.Error("OnStop callback was not called")
	}
}

// Test Server Lifecycle ===========================================================================================

// TestStartStop tests starting and stopping the server
func TestStartStop(t *testing.T) {
	srv := New().
		Port(9879).
		Build()

	// Server should not be running initially
	if srv.IsRunning() {
		t.Error("Server should not be running initially")
	}

	// Start the server
	err := srv.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Server should be running now
	if !srv.IsRunning() {
		t.Error("Server should be running after Start()")
	}

	// Stop the server
	err = srv.Stop()
	if err != nil {
		t.Fatalf("Failed to stop server: %v", err)
	}

	// Server should not be running anymore
	if srv.IsRunning() {
		t.Error("Server should not be running after Stop()")
	}
}

// TestStartTwice tests that starting twice returns an error
func TestStartTwice(t *testing.T) {
	srv := New().
		Port(9880).
		Build()

	err := srv.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	time.Sleep(10 * time.Millisecond)

	// Try to start again
	err = srv.Start()
	if err == nil {
		t.Error("Expected error when starting server twice")
	}
}

// TestStopNotRunning tests stopping a server that isn't running
func TestStopNotRunning(t *testing.T) {
	srv := New().Port(9881).Build()

	err := srv.Stop()
	if err == nil {
		t.Error("Expected error when stopping a server that isn't running")
	}
}

// TestRestart tests restarting the server
func TestRestart(t *testing.T) {
	srv := New().
		Port(9882).
		Build()

	srv.Start()
	time.Sleep(10 * time.Millisecond)

	if !srv.IsRunning() {
		t.Error("Server should be running before restart")
	}

	err := srv.Restart()
	if err != nil {
		t.Fatalf("Failed to restart server: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	if !srv.IsRunning() {
		t.Error("Server should be running after restart")
	}

	srv.Stop()
}

// Test Error Handling =================================================================================================

// TestHandlerReturnsError tests that handler errors are properly handled
func TestHandlerReturnsError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return ErrNotFound
	}

	srv := New().
		Route("/error", handler).
		Build()

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	srv.mux.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestHandlerReturnsGenericError tests generic error handling
func TestHandlerReturnsGenericError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("generic error")
	}

	srv := New().
		Route("/error", handler).
		Build()

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	srv.mux.ServeHTTP(w, req)

	// Generic errors should become 500
	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Test Getters ========================================================================================================

// TestGetPort tests GetPort()
func TestGetPort(t *testing.T) {
	srv := New().Port(3000).Build()

	if srv.GetPort() != 3000 {
		t.Errorf("Expected port 3000, got %d", srv.GetPort())
	}
}

// TestGetHost tests GetHost()
func TestGetHost(t *testing.T) {
	srv := New().Host("127.0.0.1").Build()

	if srv.GetHost() != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", srv.GetHost())
	}
}

// Test Build ==========================================================================================================

// TestBuildWithoutRoutes tests building a server with no routes
func TestBuildWithoutRoutes(t *testing.T) {
	srv := New().Port(9883).Build()

	if srv.mux == nil {
		t.Error("Expected mux to be initialized after Build()")
	}

	if srv.httpServer == nil {
		t.Error("Expected httpServer to be initialized after Build()")
	}
}

// TestAutoBuildOnStart tests that Start() calls Build() if needed
func TestAutoBuildOnStart(t *testing.T) {
	srv := New().Port(9884)

	// Don't call Build()
	err := srv.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	time.Sleep(10 * time.Millisecond)

	// Should have auto-built
	if srv.httpServer == nil {
		t.Error("Expected httpServer to be initialized by Start()")
	}
}

// Benchmark Tests =====================================================================================================

// BenchmarkServerCreation benchmarks creating a new server
func BenchmarkServerCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

// BenchmarkAddRoute benchmarks adding routes
func BenchmarkAddRoute(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		srv := New()
		srv.Route("/test", handler)
	}
}

// BenchmarkBuild benchmarks building a server
func BenchmarkBuild(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	for i := 0; i < b.N; i++ {
		srv := New().
			Route("/test1", handler).
			Route("/test2", handler).
			Route("/test3", handler)
		srv.Build()
	}
}
