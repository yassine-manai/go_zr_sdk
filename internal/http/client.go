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
	"github.com/yassine-manai/go_zr_sdk/errors"
	"github.com/yassine-manai/go_zr_sdk/logger"
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
	c.logger.Debug("making HTTP request",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
	)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed",
			logger.String("method", req.Method),
			logger.String("url", req.URL.String()),
			logger.Error(err),
		)
		return nil, errors.NewNetworkError("HTTP request failed", err)
	}

	// Log response
	c.logger.Debug("received HTTP response",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
		logger.Int("status_code", resp.StatusCode),
	)

	return resp, nil
}

// DoXML executes an XML request
func (c *Client) DoXML(ctx context.Context, method, path string, body interface{}, result interface{}) error {
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

func (c *Client) buildXMLRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
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
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("failed to read response body", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, body)
	}

	// Unmarshal XML if result is provided
	if result != nil && len(body) > 0 {
		if err := xml.Unmarshal(body, result); err != nil {
			return errors.NewSDKError(
				errors.ErrorTypeInternal,
				"failed to parse XML response",
				err,
			)
		}
	}

	return nil
}

// DoJSON executes a JSON request and unmarshals response
func (c *Client) DoJSON(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build request
	req, err := c.buildJSONRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	// Execute request
	resp, err := c.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle response
	return c.handleJSONResponse(resp, result)
}

// buildJSONRequest creates an HTTP request with JSON body
func (c *Client) buildJSONRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.config.UI.Username + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.NewSDKError(
				errors.ErrorTypeValidation,
				"failed to marshal request body",
				err,
			)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.NewNetworkError("failed to create request", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// addDefaultHeaders adds common headers to request
func (c *Client) addDefaultHeaders(req *http.Request) {
	// Basic Authentication (encode username:password)
	auth := c.config.UI.Username + ":" + c.config.UI.Password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	// Standard headers
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "thirdparty-sdk/1.0")

	// Request ID if in context
	if requestID := req.Context().Value("request_id"); requestID != nil {
		req.Header.Set("X-Request-ID", fmt.Sprintf("%v", requestID))
	}
}

// handleJSONResponse processes HTTP response and unmarshals JSON
func (c *Client) handleJSONResponse(resp *http.Response, result interface{}) error {
	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("failed to read response body", err)
	}

	// Check status code
	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, body)
	}

	// Unmarshal response
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			c.logger.Error("failed to unmarshal response",
				logger.Error(err),
				logger.String("body", string(body)),
			)
			return errors.NewSDKError(
				errors.ErrorTypeInternal,
				"failed to parse response",
				err,
			)
		}
	}

	return nil
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
		return errors.NewSDKError(
			errors.ErrorTypeInternal,
			message,
			nil,
		).WithStatusCode(statusCode)
	}
}
