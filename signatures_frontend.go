package sap

import (
	"regexp"
	"strings"
)

// Frontend framework signatures.
var frontendSignatures = []*Signature{
	// React
	{
		Framework: "React",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				Selector("div#root"),
				Selector("[data-reactroot]"),
				Selector("[data-reactid]"),
				HTMLSubstring("_reactRootContainer"),
				HTMLSubstring("__REACT_DEVTOOLS_GLOBAL_HOOK__"),
				HTMLSubstring("window.__REACT_QUERY_STATE__"),
				ScriptSrcContains("/static/js/main."),
				JSGlobal("__REACT_QUERY_STATE__"),
				JSGlobal("__REACT_DEVTOOLS_GLOBAL_HOOK__"),
			),
		},
		Version: extractReactVersion,
	},

	// Preact
	{
		Framework: "Preact",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("_preactRootContainer"),
				ScriptSrcContains("/preact"),
				HTMLSubstring("__PREACT_DEVTOOLS__"),
				JSGlobal("__PREACT_DEVTOOLS__"),
			),
		},
	},

	// Vue
	{
		Framework: "Vue",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				AttrExists("data-v"),
				Selector("[data-server-rendered=true]"),
				HTMLSubstring("__vue__"),
				HTMLSubstring("window.__VUE__"),
				HTMLSubstring("Vue.config"),
				JSGlobal("__VUE__"),
			),
		},
		Version: extractVueVersion,
	},

	// Angular
	{
		Framework: "Angular",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				SelectorAttr("*", "ng-version", ""),
				Selector("app-root"),
				AttrExists("_ngcontent"),
				AttrExists("_nghost"),
				ScriptSrcContains("zone.js"),
				ScriptSrcContains("polyfills."),
				ScriptSrcContains("runtime."),
			),
		},
		Version: extractAngularVersion,
	},

	// AngularJS
	{
		Framework: "AngularJS",
		Category:  CategorySPAFramework,
		Weight:    0.5,
		Matchers: []Matcher{
			Any(
				AttrExists("ng-app"),
				AttrExists("ng-controller"),
				ScriptSrcContains("angular.min.js"),
				ScriptSrcContains("angular.js"),
			),
		},
	},

	// Svelte
	{
		Framework: "Svelte",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("__svelte"),
				HTMLRegex(regexp.MustCompile(`class="svelte-[a-z0-9]{6}"`)),
				JSGlobal("__svelte"),
			),
		},
	},

	// Solid
	{
		Framework: "Solid",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("_$HY"),
				ScriptSrcContains("solid-js"),
				JSGlobal("_$HY"),
			),
		},
	},

	// Qwik
	{
		Framework: "Qwik",
		Category:  CategorySPAFramework,
		Weight:    0.7,
		Matchers: []Matcher{
			Any(
				AttrExists("q:container"),
				AttrExists("q:base"),
				AttrExists("q:manifest-hash"),
				Selector("script#qwikloader"),
			),
		},
	},

	// Ember
	{
		Framework: "Ember",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				Selector("div#ember-app"),
				Selector(".ember-application"),
				ScriptSrcContains("/assets/vendor-"),
			),
		},
	},

	// Lit / Web Components
	{
		Framework: "Lit",
		Category:  CategorySPAFramework,
		Weight:    0.5,
		Matchers: []Matcher{
			Any(
				ScriptSrcContains("@lit"),
				ScriptSrcContains("lit-element"),
				ScriptSrcContains("lit-html"),
				Selector("script[type=\"module\"]"),
			),
		},
	},

	// Alpine
	{
		Framework: "Alpine",
		Category:  CategorySPAFramework,
		Weight:    0.5,
		Matchers: []Matcher{
			Any(
				AttrExists("x-data"),
				AttrExists("x-bind"),
				AttrExists("x-show"),
				AttrExists("x-cloak"),
				ScriptSrcContains("alpinejs"),
			),
		},
	},

	// Stencil
	{
		Framework: "Stencil",
		Category:  CategorySPAFramework,
		Weight:    0.6,
		Matchers: []Matcher{
			Any(
				AttrExists("data-stencil-namespace"),
				ScriptSrcContains("/build/p-"),
			),
		},
	},

	// Mithril
	{
		Framework: "Mithril",
		Category:  CategorySPAFramework,
		Weight:    0.5,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("window.m"),
				ScriptSrcContains("mithril"),
			),
		},
	},

	// Backbone
	{
		Framework: "Backbone",
		Category:  CategorySPAFramework,
		Weight:    0.5,
		Matchers: []Matcher{
			Any(
				HTMLSubstring("window.Backbone"),
				ScriptSrcContains("backbone"),
			),
		},
	},

	// HTMX (not a SPA, but interactive)
	{
		Framework: "HTMX",
		Category:  CategorySPAFramework,
		Weight:    0.4,
		Matchers: []Matcher{
			Any(
				AttrExists("hx-get"),
				AttrExists("hx-post"),
				AttrExists("hx-put"),
				AttrExists("hx-delete"),
				AttrExists("hx-patch"),
				AttrExists("hx-swap"),
				AttrExists("hx-target"),
			),
		},
	},
}

// Version extractors.

func extractReactVersion(scan *Scan) string {
	doc := scan.Doc()
	if doc == nil {
		return ""
	}
	// Check meta[name="react-version"]
	if v, ok := doc.Find("meta[name='react-version']").Attr("content"); ok {
		return v
	}
	return ""
}

func extractVueVersion(scan *Scan) string {
	// Try to extract from Vue.version or __VUE__ global.
	if strings.Contains(scan.BodyStr, "Vue.version") {
		// Simple regex extraction could go here.
	}
	return ""
}

func extractAngularVersion(scan *Scan) string {
	doc := scan.Doc()
	if doc == nil {
		return ""
	}
	// Check ng-version attribute
	sel := doc.Find("[ng-version]")
	if sel.Length() > 0 {
		if v, ok := sel.Attr("ng-version"); ok {
			return v
		}
	}
	return ""
}
