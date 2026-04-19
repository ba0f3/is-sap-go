package sap

import "regexp"

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
