package sap

import "net/http"

// Category classifies what kind of technology was detected.
type Category int

const (
	CategorySPAFramework Category = iota // React, Vue, Angular, Svelte, Preact, Solid, Qwik, Ember, Lit, …
	CategoryMetaFramework               // Next.js, Nuxt, Remix, Gatsby, SvelteKit, Astro, Blazor, SolidStart, Fresh, Inertia
	CategoryHosting                     // Vercel, Netlify, Cloudflare, GitHub Pages, Fly, Render
	CategoryBundler                     // Webpack, Vite, Turbopack, esbuild, Parcel
)

func (c Category) String() string {
	switch c {
	case CategorySPAFramework:
		return "spa_framework"
	case CategoryMetaFramework:
		return "meta_framework"
	case CategoryHosting:
		return "hosting"
	case CategoryBundler:
		return "bundler"
	default:
		return "unknown"
	}
}

func (c Category) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// Framework holds detection results for a single matched technology.
type Framework struct {
	Name       string   `json:"name"`
	Category   Category `json:"category"`
	Confidence float64  `json:"confidence"`
	Version    string   `json:"version,omitempty"`
	Signals    []string `json:"signals,omitempty"`
}

// Result is the top-level output of a detection run.
type Result struct {
	URL        string            `json:"url,omitempty"`
	StatusCode int               `json:"status_code,omitempty"`
	IsSPA      bool              `json:"is_spa"`
	Confidence float64           `json:"confidence"`
	Frameworks []Framework       `json:"frameworks,omitempty"`
	Hosting    []Framework       `json:"hosting,omitempty"`
	RawHeaders http.Header       `json:"-"`
	Extras     map[string]string `json:"extras,omitempty"`
}
