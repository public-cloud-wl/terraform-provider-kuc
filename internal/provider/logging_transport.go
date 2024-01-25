package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type loggingTransport struct {
	transport http.RoundTripper
}

// RoundTrip logs the HTTP requests and responses.
func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Log the request
	tflog.Debug(ctx, "HTTP request", map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	// Perform the request using the underlying transport
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		// Log the error if the request failed
		tflog.Error(ctx, "HTTP request failed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Log the response
	tflog.Debug(ctx, "HTTP response", map[string]interface{}{
		"status": resp.Status,
	})

	return resp, nil
}

// NewLoggingTransport wraps an http.RoundTripper with logging using tflog.
func NewLoggingTransport(baseTransport http.RoundTripper) http.RoundTripper {
	return &loggingTransport{
		transport: baseTransport,
	}
}
