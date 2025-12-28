package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/yassine-manai/go_zr_sdk/config"
	"github.com/yassine-manai/go_zr_sdk/internal/errors"
	"github.com/yassine-manai/go_zr_sdk/internal/logger"
)

// Client wraps http.Client with additional functionality
type Client struct {
	httpClient *http.Client
	config     *config.Config
	logger     logger.Logger
}

// NewClient creates a new HTTP client wrapper
func NewClient(httpClient *http.Client, cfg *config.Config, log logger.Logger) *Client {
	return &Client{
		httpClient: httpClient,
		config:     cfg,
		logger:     log,
	}
}

// DoRequest executes an HTTP request and handles response
func (c *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Add default headers
	c.addDefaultHeaders(req)

	// Log request
	c.logger.Debug("making HTTP request", logger.String("method", req.Method), logger.String("url", req.URL.String()))

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed", logger.String("method", req.Method), logger.String("url", req.URL.String()), logger.Error(err))
		return nil, errors.NewNetworkError("HTTP request failed", err)
	}

	// Log response
	c.logger.Debug("received HTTP response", logger.String("method", req.Method), logger.String("url", req.URL.String()), logger.Int("status_code", resp.StatusCode))

	return resp, nil
}

// DoXML executes an XML request
func (c *Client) DoXMLRequest(ctx context.Context, method, path string, body any, result any) error {
	req, err := c.buildXMLRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	resp, err := c.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleXMLResponse(resp, result)
}

func (c *Client) buildXMLRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	url := c.config.UI.Host + c.config.UI.BasePath + path

	var bodyReader io.Reader
	if body != nil {
		xmlData, err := xml.MarshalIndent(body, "", "  ")
		if err != nil {
			return nil, err
		}
		// Add XML declaration
		fullXML := []byte(xml.Header + string(xmlData))
		bodyReader = bytes.NewReader(fullXML)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/xml")
	}

	return req, nil
}

// handleXMLResponse processes HTTP response and unmarshals XML
func (c *Client) handleXMLResponse(resp *http.Response, result interface{}) error {
	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("failed to read response body", err)
	}

	c.logger.Debug("received XML response",
		logger.Int("status_code", resp.StatusCode),
		logger.String("body", string(body)),
	)

	// Check status code
	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, body)
	}

	// Handle empty response (e.g., DELETE returns 200 with no body)
	if len(body) == 0 || result == nil {
		c.logger.Debug("empty response body, skipping unmarshal")
		return nil
	}

	// Unmarshal response
	if err := xml.Unmarshal(body, result); err != nil {
		c.logger.Error("failed to unmarshal XML response",
			logger.Error(err),
			logger.String("body", string(body)),
		)
		return errors.NewSDKError(
			errors.ErrorTypeInternal,
			"failed to parse XML response",
			err,
		)
	}

	return nil
}

// addDefaultHeaders adds common headers to request
func (c *Client) addDefaultHeaders(req *http.Request) {
	// Basic Authentication (encode username:password)
	auth := c.config.UI.Username + ":" + c.config.UI.Password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	// Standard headers
	req.Header.Set("Accept", "application/xml")

	// Request ID if in context
	if requestID := req.Context().Value("request_id"); requestID != nil {
		req.Header.Set("X-Request-ID", fmt.Sprintf("%v", requestID))
	}
}

// handleErrorResponse converts HTTP error to SDK error
func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	// Try to parse error response
	var apiErr struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}

	_ = json.Unmarshal(body, &apiErr)

	message := apiErr.Error.Message
	if message == "" {
		message = string(body)
	}

	// Map status code to error type
	switch statusCode {
	case http.StatusBadRequest:
		return errors.NewSDKError(
			errors.ErrorTypeValidation,
			message,
			nil,
		).WithStatusCode(statusCode)

	case http.StatusUnauthorized:
		return errors.NewAuthenticationError(message, nil)

	case http.StatusForbidden:
		return errors.NewAuthorizationError(message, "")

	case http.StatusNotFound:
		return errors.NewNotFoundError(message, "", "")

	case http.StatusTooManyRequests:
		return errors.NewRateLimitError(message, 0, 0, 0)

	case http.StatusServiceUnavailable:
		return errors.NewServiceUnavailableError(message, "ui-service", 0)

	default:
		return errors.NewSDKError(errors.ErrorTypeInternal, message, nil).WithStatusCode(statusCode)
	}
}
