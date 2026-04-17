package sap

// Hosting platform signatures.
var hostingSignatures = []*Signature{
	// Vercel
	{
		Framework: "Vercel",
		Category:  CategoryHosting,
		Weight:    0.8,
		Matchers: []Matcher{
			Any(
				AnyHeader("x-vercel-cache"),
				AnyHeader("x-vercel-id"),
				HeaderContains("Server", "Vercel"),
			),
		},
	},

	// Netlify
	{
		Framework: "Netlify",
		Category:  CategoryHosting,
		Weight:    0.8,
		Matchers: []Matcher{
			Any(
				HeaderContains("x-nf-request-id", ""),
				HeaderContains("Server", "Netlify"),
				HeaderContains("X-Frame-Options", "Netlify"),
			),
		},
	},

	// Cloudflare Pages
	{
		Framework: "Cloudflare Pages",
		Category:  CategoryHosting,
		Weight:    0.8,
		Matchers: []Matcher{
			Any(
				HeaderContains("cf-ray", ""),
				HeaderContains("x-served-by", "cloudflare-pages"),
			),
		},
	},

	// GitHub Pages
	{
		Framework: "GitHub Pages",
		Category:  CategoryHosting,
		Weight:    0.7,
		Matchers: []Matcher{
			Any(
				HeaderContains("Server", "GitHub.com"),
				HeaderContains("x-github-request-id", ""),
			),
		},
	},

	// Fly.io
	{
		Framework: "Fly.io",
		Category:  CategoryHosting,
		Weight:    0.7,
		Matchers: []Matcher{
			Any(
				HeaderContains("fly-request-id", ""),
				HeaderContains("Server", "Fly"),
			),
		},
	},

	// Render
	{
		Framework: "Render",
		Category:  CategoryHosting,
		Weight:    0.7,
		Matchers: []Matcher{
			Any(
				HeaderContains("x-render", ""),
				HeaderContains("Server", "Render"),
			),
		},
	},
}
