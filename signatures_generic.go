package sap

import "regexp"

var genericSignatures = []*Signature{
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
