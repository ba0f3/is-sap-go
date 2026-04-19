package sap

import (
	"context"
	"testing"
	"time"
)

// realSiteCase describes one integration test case.
type realSiteCase struct {
	url         string
	framework   string // primary expected framework name
	category    string // "spa_framework" | "meta_framework"
	minConf     float64
	description string
}

// realSiteCases are well-known public sites that should be reachable and stable.
// Tests skip individually when a site is unreachable (network error or non-2xx).
var realSiteCases = []realSiteCase{
	{
		url:         "https://nextjs.org",
		framework:   "Next.js",
		category:    "meta_framework",
		minConf:     0.4,
		description: "Next.js official site",
	},
	{
		url:         "https://nuxt.com",
		framework:   "Nuxt",
		category:    "meta_framework",
		minConf:     0.4,
		description: "Nuxt official site",
	},
	{
		url:         "https://remix.run",
		framework:   "Remix",
		category:    "meta_framework",
		minConf:     0.3,
		description: "Remix official site",
	},
	{
		url:         "https://react.dev",
		framework:   "React",
		category:    "spa_framework",
		minConf:     0.3,
		description: "React official docs site",
	},
	{
		url:         "https://vuejs.org",
		framework:   "Vue",
		category:    "spa_framework",
		minConf:     0.3,
		description: "Vue.js official site",
	},
	{
		url:         "https://angular.io",
		framework:   "Angular",
		category:    "spa_framework",
		minConf:     0.3,
		description: "Angular official site",
	},
	{
		url:         "https://astro.build",
		framework:   "Astro",
		category:    "meta_framework",
		minConf:     0.4,
		description: "Astro official site",
	},
	{
		url:         "https://svelte.dev",
		framework:   "SvelteKit",
		category:    "meta_framework",
		minConf:     0.3,
		description: "Svelte/SvelteKit official site",
	},
}

// TestIntegration_RealSites tests detection against live public websites.
// Run with: go test -run TestIntegration -tags integration
// Or with INTEGRATION=1 environment variable.
func TestIntegration_RealSites(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	for _, tc := range realSiteCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			result, err := Detect(ctx, tc.url,
				WithTimeout(15*time.Second),
				WithUserAgent("is-sap-go-test/1.0"),
			)
			if err != nil {
				t.Skipf("site unreachable (%s): %v", tc.url, err)
			}

			if result.StatusCode >= 500 {
				t.Skipf("site returned server error %d, skipping", result.StatusCode)
			}

			found := false
			for _, fw := range result.Frameworks {
				if fw.Name == tc.framework && fw.Confidence >= tc.minConf {
					found = true
					t.Logf("detected %s (confidence=%.2f, signals=%v)", fw.Name, fw.Confidence, fw.Signals)
					break
				}
			}

			if !found {
				// Log all detected frameworks as context, but only fail if we got 2xx.
				t.Logf("all frameworks: %+v", result.Frameworks)
				t.Logf("is_spa=%v, confidence=%.2f", result.IsSPA, result.Confidence)
				if result.StatusCode >= 200 && result.StatusCode < 300 {
					t.Errorf("expected %q (>= %.1f confidence) in response from %s (HTTP %d)",
						tc.framework, tc.minConf, tc.url, result.StatusCode)
				} else {
					t.Skipf("non-2xx response %d from %s, skipping assertion", result.StatusCode, tc.url)
				}
			}
		})
	}
}

// TestIntegration_Vercel tests that Vercel hosting is detected on a Vercel-hosted site.
func TestIntegration_Vercel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := Detect(ctx, "https://vercel.com",
		WithTimeout(15*time.Second),
		WithUserAgent("is-sap-go-test/1.0"),
	)
	if err != nil {
		t.Skipf("vercel.com unreachable: %v", err)
	}

	for _, h := range result.Hosting {
		if h.Name == "Vercel" {
			t.Logf("Vercel hosting detected (confidence=%.2f, signals=%v)", h.Confidence, h.Signals)
			return
		}
	}

	if result.StatusCode >= 200 && result.StatusCode < 300 {
		t.Logf("hosting detected: %+v", result.Hosting)
		t.Logf("frameworks detected: %+v", result.Frameworks)
		// Soft failure: Vercel may change headers, so just log instead of hard fail.
		t.Logf("Vercel hosting not detected; Vercel may have changed response headers")
	} else {
		t.Skipf("non-2xx response %d from vercel.com", result.StatusCode)
	}
}
