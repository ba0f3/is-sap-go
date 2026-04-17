package sap

// Meta-framework signatures.
var metaFrameworkSignatures = []*Signature{
	// Next.js
	{
		Framework: "Next.js",
		Category:  CategoryMetaFramework,
		Weight:    0.95,
		Matchers: []Matcher{
			Any(
				Selector("div#__next"),
				Selector("script#__NEXT_DATA__"),
				LinkHrefContains("/_next/static/"),
				ScriptSrcContains("/_next/"),
				HeaderContains("x-nextjs-cache", ""),
				HeaderContains("x-nextjs-prerender", ""),
				HeaderContains("x-nextjs-matched-path", ""),
				CookieNameHas("next-auth.session-token"),
				JSGlobal("__NEXT_DATA__"),
			),
		},
		Implies: []string{"React"},
	},

	// Nuxt
	{
		Framework: "Nuxt",
		Category:  CategoryMetaFramework,
		Weight:    0.95,
		Matchers: []Matcher{
			Any(
				Selector("div#__nuxt"),
				Selector("div#__layout"),
				LinkHrefContains("/_nuxt/"),
				ScriptSrcContains("/_nuxt/"),
				Selector("script#__NUXT_DATA__"),
				JSGlobal("__NUXT__"),
			),
		},
		Implies: []string{"Vue"},
	},

	// Remix
	{
		Framework: "Remix",
		Category:  CategoryMetaFramework,
		Weight:    0.95,
		Matchers: []Matcher{
			Any(
				JSGlobal("__remixContext"),
				JSGlobal("__remixManifest"),
				JSGlobal("__remixRouteModules"),
				HTMLSubstring("<!-- remix -->"),
			),
		},
		Implies: []string{"React"},
	},

	// Gatsby
	{
		Framework: "Gatsby",
		Category:  CategoryMetaFramework,
		Weight:    0.9,
		Matchers: []Matcher{
			Any(
				Selector("div#___gatsby"),
				HTMLSubstring("window.___gatsby"),
				LinkHrefContains("/page-data/"),
				HTMLSubstring("gatsby-chunk-mapping"),
			),
		},
		Implies: []string{"React"},
	},

	// SvelteKit
	{
		Framework: "SvelteKit",
		Category:  CategoryMetaFramework,
		Weight:    0.9,
		Matchers: []Matcher{
			Any(
				LinkHrefContains("/_app/immutable/"),
				AttrExists("data-sveltekit-preload-code"),
				AttrExists("data-sveltekit-preload-data"),
				HTMLSubstring("__sveltekit_"),
			),
		},
		Implies: []string{"Svelte"},
	},

	// Astro
	{
		Framework: "Astro",
		Category:  CategoryMetaFramework,
		Weight:    0.9,
		Matchers: []Matcher{
			Any(
				Selector("astro-island"),
				AttrExists("data-astro-cid"),
				LinkHrefContains("/_astro/"),
				ScriptSrcContains("/_astro/"),
			),
		},
	},

	// Blazor
	{
		Framework: "Blazor",
		Category:  CategoryMetaFramework,
		Weight:    0.95,
		Matchers: []Matcher{
			Any(
				ScriptSrcContains("_framework/blazor.webassembly.js"),
				ScriptSrcContains("_framework/blazor.server.js"),
				LinkHrefContains("/_framework/blazor.boot.json"),
			),
		},
	},

	// SolidStart
	{
		Framework: "SolidStart",
		Category:  CategoryMetaFramework,
		Weight:    0.9,
		Matchers: []Matcher{
			Any(
				LinkHrefContains("/_build/"),
				AttrExists("data-hk"),
			),
		},
		Implies: []string{"Solid"},
	},

	// Fresh (Deno)
	{
		Framework: "Fresh",
		Category:  CategoryMetaFramework,
		Weight:    0.9,
		Matchers: []Matcher{
			Any(
				ScriptSrcContains("/_frsh/"),
				HTMLSubstring("<!--frsh-"),
			),
		},
	},

	// Inertia.js
	{
		Framework: "Inertia.js",
		Category:  CategoryMetaFramework,
		Weight:    0.85,
		Matchers: []Matcher{
			Any(
				SelectorAttr("div#app", "data-page", ""),
				HeaderContains("X-Inertia", "true"),
			),
		},
	},

	// Analog (Angular + signal-based)
	{
		Framework: "Analog",
		Category:  CategoryMetaFramework,
		Weight:    0.8,
		Matchers: []Matcher{
			Any(
				ScriptSrcContains("_analog_"),
			),
		},
		Implies: []string{"Angular"},
	},
}
