package sap

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Options configures detection behavior for [Detect], [DetectFromResponse],
// and [DetectFromHTML].
type Options struct {
	Timeout         time.Duration
	UserAgent       string
	MaxBodyBytes    int64
	FollowRedirects bool
	EnableJSEval    bool
	HTTPClient      *http.Client
}

// Option is a functional option for configuring [Options].
type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Timeout:         10 * time.Second,
		UserAgent:       "is-sap-go/1.0",
		MaxBodyBytes:    2 * 1024 * 1024, // 2 MiB
		FollowRedirects: true,
		EnableJSEval:    false,
	}
}

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

// WithEnableJSEval enables the optional goja JS sandbox for extracting
// window.__* globals from inline scripts.
func WithEnableJSEval(enable bool) Option {
	return func(o *Options) {
		o.EnableJSEval = enable
	}
}

// WithHTTPClient sets a custom HTTP client for [Detect].
func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = c
	}
}

// Detect fetches rawURL and detects the SPA framework used by the page.
func Detect(ctx context.Context, rawURL string, opts ...Option) (*Result, error) {
	options := defaultOptions()
	applyOptions(options, opts...)

	resp, body, err := fetch(ctx, rawURL, options)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	result, err := DetectFromHTML(string(body), resp.Header, opts...)
	if err != nil {
		return result, err
	}
	result.StatusCode = resp.StatusCode

	if u, parseErr := parseURL(rawURL); parseErr == nil {
		result.URL = u
	}

	return result, nil
}

// DetectFromResponse analyzes a pre-fetched HTTP response.
func DetectFromResponse(resp *http.Response, body []byte, opts ...Option) (*Result, error) {
	result, err := DetectFromHTML(string(body), resp.Header, opts...)
	if err != nil {
		return result, err
	}
	result.StatusCode = resp.StatusCode

	if u, parseErr := parseURL(resp.Request.URL.String()); parseErr == nil {
		result.URL = u
	}

	return result, nil
}

// DetectFromHTML analyzes raw HTML and headers to detect SPA frameworks.
func DetectFromHTML(html string, headers http.Header, opts ...Option) (*Result, error) {
	options := defaultOptions()
	applyOptions(options, opts...)

	if headers == nil {
		headers = make(http.Header)
	}

	scan := newScan("", headers, []byte(html), options)

	result := &Result{
		IsSPA:      false,
		Confidence: 0,
		Frameworks: []Framework{},
		Hosting:    []Framework{},
		RawHeaders: headers,
		Extras:     make(map[string]string),
	}

	scoreSignatures(scan, result)

	suppressRedundant(result)

	sortFrameworks(result.Frameworks)
	sortFrameworks(result.Hosting)

	result.IsSPA = determineSPA(result)
	result.Confidence = calculateConfidence(result)

	return result, nil
}

func parseURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
