package mux

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRouter_RemovePathByName(t *testing.T) {
	router := NewRouter()

	router.HandleFuncWithName("health", "/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	ctx, closeFn := context.WithCancel(context.Background())
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go srv.ListenAndServe()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			srv.Close()
		}
	}(ctx)
	time.Sleep(time.Second)
	req, _ := http.Get("http://127.0.0.1:8000/api/health")
	data, _ := io.ReadAll(req.Body)
	if string(data) != "{\"ok\":true}\n" {
		t.Error("/api/health error")
	}
	err := router.RemovePathByName("health")
	if err != nil {
		t.Error(err)
	}
	req, _ = http.Get("http://127.0.0.1:8000/api/health")
	data, _ = io.ReadAll(req.Body)
	if string(data) != "404 page not found\n" {
		t.Error("remove /api/health error")
	}
	closeFn()
	time.Sleep(time.Second)
}
