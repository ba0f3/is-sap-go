package sap

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Scan holds all the data about a single URL inspection, with lazy parsing.
type Scan struct {
	URL       *url.URL
	Headers   http.Header
	Cookies   []*http.Cookie
	Body      []byte
	BodyStr   string
	BodyLower string
	opts      *Options

	docOnce sync.Once
	doc     *goquery.Document

	blobsOnce sync.Once
	jsBlobs   []JSBlob
}

func newScan(rawURL string, headers http.Header, body []byte, opts *Options) *Scan {
	if headers == nil {
		headers = make(http.Header)
	}
	s := &Scan{
		Headers:   headers,
		Body:      body,
		BodyStr:   string(body),
		BodyLower: strings.ToLower(string(body)),
		opts:      opts,
		Cookies:   (&http.Response{Header: headers}).Cookies(),
	}
	if rawURL != "" {
		s.URL, _ = url.Parse(rawURL)
	}
	return s
}

// Doc returns a lazily-parsed goquery document. May return nil on parse error.
func (s *Scan) Doc() *goquery.Document {
	s.docOnce.Do(func() {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(s.BodyStr))
		if err == nil {
			s.doc = doc
		}
	})
	return s.doc
}

// JSBlobs runs the goja inline-script sandbox and returns harvested window.__* globals.
// Returns nil if EnableJSEval is false.
func (s *Scan) JSBlobs() []JSBlob {
	if !s.opts.EnableJSEval {
		return nil
	}
	s.blobsOnce.Do(func() {
		s.jsBlobs = ExtractJSDataFromHTML(s.BodyStr)
	})
	return s.jsBlobs
}
