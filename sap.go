package sap

import (
	"context"
	"net/http"
	"time"
)

// Options configures detection behavior.
type Options struct {
	Timeout         time.Duration
	UserAgent       string
	MaxBodyBytes    int64
	FollowRedirects bool
	EnableJSEval    bool
	HTTPClient      *http.Client
}

// Option is a functional option for configuring Options.
type Option func(*Options)

// defaultOptions returns Options with sensible defaults.
func defaultOptions() *Options {
	return &Options{
		Timeout:         10 * time.Second,
		UserAgent:       "is-sap-go/1.0",
		MaxBodyBytes:    2 * 1024 * 1024, // 2 MiB
		FollowRedirects: true,
		EnableJSEval:    false,
	}
}

// applyOptions merges Option funcs into the base.
func applyOptions(opts *Options, fns ...Option) {
	for _, f := range fns {
		f(opts)
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(o *Options) {
		o.UserAgent = ua
	}
}

// WithMaxBodyBytes sets the maximum response body size to read.
func WithMaxBodyBytes(n int64) Option {
	return func(o *Options) {
		o.MaxBodyBytes = n
	}
}

// WithFollowRedirects controls 3xx redirect behavior.
func WithFollowRedirects(follow bool) Option {
	return func(o *Options) {
		o.FollowRedirects = follow
	}
}

// WithEnableJSEval enables the optional goja JS sandbox.
func WithEnableJSEval(enable bool) Option {
	return func(o *Options) {
		o.EnableJSEval = enable
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = c
	}
}

// Detect fetches rawURL and detects the SPA framework.
func Detect(ctx context.Context, rawURL string, opts ...Option) (*Result, error) {
	options := defaultOptions()
	applyOptions(options, opts...)

	resp, body, err := fetch(ctx, rawURL, options)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return DetectFromResponse(resp, body, opts...)
}

// DetectFromResponse analyzes a pre-fetched HTTP response.
func DetectFromResponse(resp *http.Response, body []byte, opts ...Option) (*Result, error) {
	options := defaultOptions()
	applyOptions(options, opts...)

	return DetectFromHTML(string(body), resp.Header, append(opts,
		WithTimeout(0), // ignore timeout for already-fetched content
	)...)
}

// DetectFromHTML is the lowest-level primitive; analyzes raw HTML and headers.
func DetectFromHTML(html string, headers http.Header, opts ...Option) (*Result, error) {
	options := defaultOptions()
	applyOptions(options, opts...)

	if headers == nil {
		headers = make(http.Header)
	}

	// Create a scan context.
	scan := newScan("", headers, []byte(html), options)

	// Run all signatures through the registry.
	result := &Result{
		IsSPA:      false,
		Confidence: 0,
		Frameworks: []Framework{},
		Hosting:    []Framework{},
		RawHeaders: headers,
		Extras:     make(map[string]string),
	}

	// Score all signatures.
	scoreSignatures(scan, result)

	// Filter redundant frameworks (e.g., suppress React if Next.js is high-confidence).
	suppressRedundant(result)

	// Sort by confidence descending.
	sortFrameworks(result.Frameworks)
	sortFrameworks(result.Hosting)

	// Determine if it's a SPA.
	result.IsSPA = determineSPA(result)

	// Overall confidence = max of SPA frameworks + bump from generic heuristics.
	result.Confidence = calculateConfidence(result)

	return result, nil
}
