package server

import (
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestNotGracefulShutdown(t *testing.T) {
	srv := RunNotGracefulServer()

	time.Sleep(500 * time.Millisecond)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/long", nil)
		if err != nil {
			t.Errorf("Fail to create http request : %v", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Logf("Expected request to success, but it fail: %v", err)
			return
		}
		defer resp.Body.Close()
	}()
	time.Sleep(1 * time.Second)

	if err := srv.Close(); err != nil {
		t.Errorf("Server closing error : %v", err)
	}
	wg.Wait()
}

func TestGracefulShutdown(t *testing.T) {
	shutdownChan := make(chan struct{})
	done := make(chan bool, 1)
	_, err := RunGracefulServer(done, shutdownChan)
	if err != nil {
		t.Errorf("Fail to launch server : %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/long", nil)
		if err != nil {
			t.Errorf("Fail to create http request : %v", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Fail to call http request : %v", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Fail to read http body : %v", err)
			return
		}

		if string(body) != TASK_COMPLETE {
			t.Errorf("Expect : %s, but got %s", TASK_COMPLETE, string(body))
			return
		}
	}()

	time.Sleep(1 * time.Second)
	close(shutdownChan)
	wg.Wait()
	<-done
}
