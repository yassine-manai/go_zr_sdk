package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/yassine-manai/go_zr_sdk/config"
)

// createHTTPClient creates a configured HTTP client
func createHTTPClient(cfg *config.Config) *http.Client {
	timeout := cfg.Timeout
	if cfg.UI.Timeout > 0 {
		timeout = cfg.UI.Timeout
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			// Connection pooling
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.UI.InsecureSkipVerify,
			},

			// Timeouts
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
