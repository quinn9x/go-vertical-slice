package app_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

const testServerAddress = "127.0.0.1:0"

func TestGracefulShutdown_WaitsForActiveRequest(t *testing.T) {
	t.Parallel()

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	requestFinished := make(chan struct{})

	e := echo.New()

	e.GET("/slow", func(c *echo.Context) error {
		close(requestStarted)

		<-releaseRequest

		close(requestFinished)

		return c.String(http.StatusOK, "ok")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addressCh := make(chan string, 1)
	serverErr := make(chan error, 1)

	startConfig := echo.StartConfig{
		Address:         testServerAddress,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 2 * time.Second,
		ListenerAddrFunc: func(addr net.Addr) {
			addressCh <- addr.String()
		},
	}

	go func() {
		serverErr <- startConfig.Start(ctx, e)
	}()

	address := <-addressCh

	clientDone := make(chan error, 1)

	go func() {
		clientDone <- doRequest(context.Background(), address+"/slow")
	}()

	<-requestStarted

	cancel()

	select {
	case <-requestFinished:
		t.Fatal("request finished before it was released")
	default:
	}

	close(releaseRequest)

	select {
	case <-requestFinished:
	case <-time.After(1 * time.Second):
		t.Fatal("active request did not finish")
	}

	require.NoError(t, <-clientDone)
	require.NoError(t, <-serverErr)
}

func TestGracefulShutdown_ReturnsAfterContextCancellation(t *testing.T) {
	t.Parallel()

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addressCh := make(chan string, 1)
	serverErr := make(chan error, 1)

	startConfig := echo.StartConfig{
		Address:         testServerAddress,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 2 * time.Second,
		ListenerAddrFunc: func(addr net.Addr) {
			addressCh <- addr.String()
		},
	}

	go func() {
		serverErr <- startConfig.Start(ctx, e)
	}()

	<-addressCh

	cancel()

	select {
	case err := <-serverErr:
		require.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("server did not shut down after context cancellation")
	}
}

func TestGracefulShutdown_StopsWaitingAfterTimeout(t *testing.T) {
	t.Parallel()

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})

	e := echo.New()

	e.GET("/slow", func(c *echo.Context) error {
		close(requestStarted)

		<-releaseRequest

		return c.String(http.StatusOK, "ok")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addressCh := make(chan string, 1)
	serverErr := make(chan error, 1)

	const gracefulTimeout = 100 * time.Millisecond

	startConfig := echo.StartConfig{
		Address:         testServerAddress,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: gracefulTimeout,
		ListenerAddrFunc: func(addr net.Addr) {
			addressCh <- addr.String()
		},
	}

	go func() {
		serverErr <- startConfig.Start(ctx, e)
	}()

	address := <-addressCh

	clientDone := make(chan error, 1)

	go func() {
		clientDone <- doRequest(context.Background(), address+"/slow")
	}()

	<-requestStarted

	start := time.Now()

	cancel()

	select {
	case err := <-serverErr:
		require.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("server did not stop after graceful shutdown timeout")
	}

	elapsed := time.Since(start)

	require.GreaterOrEqual(t, elapsed, gracefulTimeout)

	close(releaseRequest)

	select {
	case <-clientDone:
	case <-time.After(1 * time.Second):
		t.Fatal("blocked request did not finish after release")
	}
}

func doRequest(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+url,
		http.NoBody,
	)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}
