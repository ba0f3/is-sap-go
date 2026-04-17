package sap

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// TestDetectFromHTML_Fixtures tests detection against fixture HTML files.
func TestDetectFromHTML_Fixtures(t *testing.T) {
	tests := []struct {
		fixture       string
		expectedFW    string
		minConfidence float64
	}{
		{"nextjs-vercel.html", "Next.js", 0.5},
		{"nuxt3.html", "Nuxt", 0.5},
		{"react.html", "React", 0.4},
		{"angular.html", "Angular", 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			// Load fixture file
			path := filepath.Join("testdata", tt.fixture)
			html, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read fixture: %v", err)
			}

			// Detect
			result, err := DetectFromHTML(string(html), http.Header{})
			if err != nil {
				t.Fatalf("DetectFromHTML failed: %v", err)
			}

			// Check for expected framework
			found := false
			for _, fw := range result.Frameworks {
				if fw.Name == tt.expectedFW && fw.Confidence >= tt.minConfidence {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected framework %q with confidence >= %.1f, got: %+v",
					tt.expectedFW, tt.minConfidence, result.Frameworks)
			}

			// Should be marked as SPA
			if !result.IsSPA {
				t.Errorf("expected IsSPA=true, got false")
			}
		})
	}
}

// TestMatcherBuilders_HTMLSubstring tests the HTML substring matcher.
func TestMatcherBuilders_HTMLSubstring(t *testing.T) {
	html := `<html><body><div id="root"></div></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := HTMLSubstring("div id=\"root\"")
	hit, ev := m(scan)

	if !hit {
		t.Errorf("expected match for 'div id=\"root\"', got no match")
	}
	if ev == "" {
		t.Errorf("expected non-empty evidence string")
	}
}

// TestMatcherBuilders_Selector tests the selector matcher.
func TestMatcherBuilders_Selector(t *testing.T) {
	html := `<html><body><app-root></app-root></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := Selector("app-root")
	hit, ev := m(scan)

	if !hit {
		t.Errorf("expected match for 'app-root' selector, got no match")
	}
	if ev == "" {
		t.Errorf("expected non-empty evidence string")
	}
}

// TestMatcherBuilders_HeaderContains tests the header matcher.
func TestMatcherBuilders_HeaderContains(t *testing.T) {
	headers := http.Header{}
	headers.Set("X-Vercel-ID", "abc123")

	scan := newScan("", headers, []byte(""), defaultOptions())

	m := HeaderContains("X-Vercel-ID", "abc")
	hit, ev := m(scan)

	if !hit {
		t.Errorf("expected match for header X-Vercel-ID, got no match")
	}
	if ev == "" {
		t.Errorf("expected non-empty evidence string")
	}
}

// TestMatcherBuilders_CookieNameHas tests the cookie matcher.
func TestMatcherBuilders_CookieNameHas(t *testing.T) {
	headers := http.Header{}
	headers.Set("Set-Cookie", "next-auth.session-token=xyz123; Path=/")

	scan := newScan("", headers, []byte(""), defaultOptions())

	m := CookieNameHas("next-auth.session-token")
	hit, ev := m(scan)

	if !hit {
		t.Errorf("expected match for cookie 'next-auth.session-token', got no match")
	}
	if ev == "" {
		t.Errorf("expected non-empty evidence string")
	}
}

// TestMatcherBuilders_Any tests the OR combinator.
func TestMatcherBuilders_Any(t *testing.T) {
	html := `<html><body><div id="root"></div></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := Any(
		HTMLSubstring("nonexistent"),
		HTMLSubstring("root"),
	)
	hit, _ := m(scan)

	if !hit {
		t.Errorf("expected Any() match, got no match")
	}
}

// TestMatcherBuilders_All tests the AND combinator.
func TestMatcherBuilders_All(t *testing.T) {
	html := `<html><body><div id="root"></div></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := All(
		HTMLSubstring("root"),
		HTMLSubstring("div"),
	)
	hit, _ := m(scan)

	if !hit {
		t.Errorf("expected All() match, got no match")
	}
}

// TestDetectFromHTML_MinimalHTML tests with empty/minimal HTML.
func TestDetectFromHTML_MinimalHTML(t *testing.T) {
	html := `<html><body></body></html>`
	result, err := DetectFromHTML(html, http.Header{})
	if err != nil {
		t.Fatalf("DetectFromHTML failed: %v", err)
	}

	// Should not be marked as SPA without framework signals
	if result.IsSPA {
		t.Errorf("expected IsSPA=false for minimal HTML, got true")
	}
}

// TestScan_LazyDoc verifies lazy parsing of goquery.Document.
func TestScan_LazyDoc(t *testing.T) {
	html := `<html><body><div id="test">content</div></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	// First call should parse
	doc1 := scan.Doc()
	if doc1 == nil {
		t.Errorf("expected non-nil Document")
	}

	// Second call should return same instance (memoized)
	doc2 := scan.Doc()
	if doc1 != doc2 {
		t.Errorf("expected same Document instance (memoization), got different")
	}
}

// TestScan_JSBlobs_Disabled tests that JSBlobs returns nil when disabled.
func TestScan_JSBlobs_Disabled(t *testing.T) {
	html := `<html><body><script>window.__TEST__ = 123;</script></body></html>`
	opts := defaultOptions()
	opts.EnableJSEval = false
	scan := newScan("", http.Header{}, []byte(html), opts)

	blobs := scan.JSBlobs()
	if blobs != nil {
		t.Errorf("expected nil JSBlobs when EnableJSEval=false, got %v", blobs)
	}
}

// TestMatcherBuilders_LinkHref tests the link href matcher.
func TestMatcherBuilders_LinkHref(t *testing.T) {
	html := `<html><head><link rel="stylesheet" href="/_next/static/app.css"></head></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := LinkHrefContains("/_next/static")
	hit, _ := m(scan)

	if !hit {
		t.Errorf("expected match for link href containing '/_next/static'")
	}
}

// TestMatcherBuilders_ScriptSrc tests the script src matcher.
func TestMatcherBuilders_ScriptSrc(t *testing.T) {
	html := `<html><body><script src="/_next/static/main.js"></script></body></html>`
	scan := newScan("", http.Header{}, []byte(html), defaultOptions())

	m := ScriptSrcContains("/_next/")
	hit, _ := m(scan)

	if !hit {
		t.Errorf("expected match for script src containing '/_next/'")
	}
}
