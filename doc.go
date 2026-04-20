// Package sap detects Single Page Application (SPA) frameworks from HTML
// responses and HTTP headers.
//
// It identifies meta-frameworks (Next.js, Nuxt, Remix, etc.), SPA libraries
// (React, Vue, Angular, etc.), hosting platforms (Vercel, Netlify, etc.),
// and generic SPA heuristics (hashed bundles, modulepreload, noscript tags).
//
// The primary entry points are [Detect] for live URL detection and
// [DetectFromHTML] for pre-fetched HTML analysis.
//
// # Quick Start
//
//	result, err := sap.Detect(ctx, "https://nextjs.org")
//	if err != nil { /* handle error */ }
//	fmt.Printf("IsSPA: %v, Confidence: %.2f\n", result.IsSPA, result.Confidence)
//
// # Functional Options
//
// Use [Option] functions to configure detection behavior:
//
//	result, err := sap.Detect(ctx, url,
//	    sap.WithTimeout(5*time.Second),
//	    sap.WithEnableJSEval(true),
//	    sap.WithUserAgent("my-scanner/1.0"),
//	)
package sap
