package sap

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// fetch performs an HTTP request and returns the response and body.
func fetch(ctx context.Context, rawURL string, opts *Options) (*http.Response, []byte, error) {
	client := opts.HTTPClient
	if client == nil {
		tr := &http.Transport{
			DisableKeepAlives: true,
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   opts.Timeout,
		}
		if !opts.FollowRedirects {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	// Read body with size limit.
	limitedBody := io.LimitReader(resp.Body, opts.MaxBodyBytes)
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, body, nil
}
