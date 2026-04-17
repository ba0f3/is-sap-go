package sap

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Generic SPA heuristics (weak signals used when no framework is detected).
var genericSignatures = []*Signature{
	// Generic root app div pattern
	{
		Framework: "SPA (generic)",
		Category:  CategorySPAFramework,
		Weight:    0.4,
		Matchers: []Matcher{
			Any(
				Selector("div#app"),
				Selector("div#root"),
				Selector("div#main"),
				Selector("div#__nuxt"),
				Selector("div#__next"),
				Selector("div#___gatsby"),
			),
		},
	},

	// JavaScript required message (noscript tag)
	{
		Framework: "SPA (noscript)",
		Category:  CategorySPAFramework,
		Weight:    0.3,
		Matchers: []Matcher{
			Any(
				Selector("noscript"),
			),
		},
	},

	// Hashed bundle filenames
	{
		Framework: "SPA (hashed-bundles)",
		Category:  CategorySPAFramework,
		Weight:    0.35,
		Matchers: []Matcher{
			Any(
				HTMLRegex(regexp.MustCompile(`[a-z0-9]+\.[a-f0-9]{8,}\.(js|css|mjs)`)),
			),
		},
	},

	// modulepreload link tags
	{
		Framework: "SPA (modulepreload)",
		Category:  CategorySPAFramework,
		Weight:    0.3,
		Matchers: []Matcher{
			Any(
				LinkHrefContains("rel=\"modulepreload\""),
			),
		},
	},

	// Service worker registration
	{
		Framework: "SPA (service-worker)",
		Category:  CategorySPAFramework,
		Weight:    0.4,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("navigator.serviceWorker.register"),
			),
		},
	},

	// PWA manifest
	{
		Framework: "SPA (manifest)",
		Category:  CategorySPAFramework,
		Weight:    0.25,
		Matchers: []Matcher{
			Any(
				LinkHrefContains("rel=\"manifest\""),
			),
		},
	},

	// Hash-based routing
	{
		Framework: "SPA (hash-routing)",
		Category:  CategorySPAFramework,
		Weight:    0.3,
		Matchers: []Matcher{
			Any(
				Selector("base[href*='#/']"),
			),
		},
	},
}

// Generic heuristics integrated into scoring.
// These are already handled by the genericSignatures above, but we can add
// more sophisticated detection in the future (e.g., analyzing bundle structure).
func detectGenericSPA(scan *Scan) bool {
	// Count hashed bundles in script/link tags.
	doc := scan.Doc()
	if doc == nil {
		return false
	}

	hashedBundleRe := regexp.MustCompile(`\.[a-f0-9]{8,}\.(js|mjs|css)$`)

	hashedScripts := 0
	doc.Find("script[src]").Each(func(_ int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && hashedBundleRe.MatchString(src) {
			hashedScripts++
		}
	})

	hashedLinks := 0
	doc.Find("link[href][rel~='stylesheet']").Each(func(_ int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok && hashedBundleRe.MatchString(href) {
			hashedLinks++
		}
	})

	// If >= 3 hashed modulepreload links, strong SPA signal.
	modulepreloadCount := 0
	doc.Find("link[rel~='modulepreload']").Each(func(_ int, s *goquery.Selection) {
		modulepreloadCount++
	})

	if modulepreloadCount >= 3 {
		return true
	}

	// If body has minimal text and a single app root div, likely SPA.
	bodyText := strings.TrimSpace(doc.Find("body").Text())
	if len(bodyText) < 1000 && countAppRoots(doc) == 1 {
		return true
	}

	return false
}

// countAppRoots counts div elements with IDs in the known app root list.
func countAppRoots(doc *goquery.Document) int {
	appIds := []string{"app", "root", "main", "__nuxt", "__next", "___gatsby"}
	count := 0
	for _, id := range appIds {
		if doc.Find("div#" + id).Length() > 0 {
			count++
		}
	}
	return count
}
