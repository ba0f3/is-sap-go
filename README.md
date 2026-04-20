<div align="center">

# is-sap-go

**Detect Single Page Application frameworks from HTML responses**

[![CI](https://github.com/ba0f3/is-sap-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ba0f3/is-sap-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go/dev/badge/github.com/ba0f3/is-sap-go.svg)](https://pkg.go/dev/github.com/ba0f3/is-sap-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ba0f3/is-sap-go)](https://goreportcard.com/report/github.com/ba0f3/is-sap-go)

</div>

---

`is-sap-go` inspects an HTTP response (or raw HTML) and determines whether a page is a **Single Page Application**, identifying the specific framework, meta-framework, and hosting platform with confidence scores.

## Features

- **30+ signatures** covering React, Vue, Angular, Svelte, Next.js, Nuxt, Remix, Astro, SvelteKit, and more
- **Confidence-weighted scoring** with hit multiplier for multiple matching signals
- **Redundancy suppression** — detects Next.js without also reporting React
- **Hosting detection** — Vercel, Netlify, Cloudflare Pages, GitHub Pages, Fly.io, Render
- **Optional JS sandbox** via [goja](https://github.com/dop251/goja) to extract `window.__*` globals from inline scripts
- **Zero HTTP dependencies** for the detection path — bring your own `*http.Client`
- **Functional options** API for clean configuration

## Install

```bash
go get github.com/ba0f3/is-sap-go
```

**Requires Go 1.22+**

## Quick Start

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    sap "github.com/ba0f3/is-sap-go"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := sap.Detect(ctx, "https://nextjs.org",
        sap.WithTimeout(10*time.Second),
        sap.WithUserAgent("my-scanner/1.0"),
    )
    if err != nil {
        panic(err)
    }

    data, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(data))
}
```

Output:

```json
{
  "url": "https://nextjs.org",
  "status_code": 200,
  "is_spa": true,
  "confidence": 0.95,
  "frameworks": [
    {
      "name": "Next.js",
      "category": "meta_framework",
      "confidence": 0.95,
      "signals": ["selector \"div#__next\"", "script src contains \"/_next/\""]
    }
  ],
  "hosting": [
    {
      "name": "Vercel",
      "category": "hosting",
      "confidence": 0.80,
      "signals": ["header x-vercel-id: abc123"]
    }
  ],
  "headers": {
    "Content-Type": ["text/html; charset=utf-8"],
    "X-Vercel-Id": ["abc123"]
  }
}
```

## API

### `Detect(ctx, url, ...Option) (*Result, error)`

Fetches the URL and runs all detection signatures.

### `DetectFromResponse(resp, body, ...Option) (*Result, error)`

Detects from a pre-fetched `*http.Response` and its body bytes.

### `DetectFromHTML(html, headers, ...Option) (*Result, error)`

Lowest-level entry point. Analyzes raw HTML string and `http.Header`.

### Options

| Option | Type | Default | Description |
|---|---|---|---|
| `WithTimeout` | `time.Duration` | `10s` | HTTP request timeout |
| `WithUserAgent` | `string` | `"is-sap-go/1.0"` | User-Agent header |
| `WithMaxBodyBytes` | `int64` | `2 MiB` | Max response body to read |
| `WithFollowRedirects` | `bool` | `true` | Follow HTTP 3xx redirects |
| `WithEnableJSEval` | `bool` | `false` | Enable goja JS sandbox |
| `WithHTTPClient` | `*http.Client` | *(new)* | Custom HTTP client |

## Detection Methodology

Each **signature** defines a framework name, category, weight, and a list of **matchers**:

| Weight | Meaning |
|---|---|
| `0.95` | Definitive signal (e.g., `#__next` selector for Next.js) |
| `0.6–0.8` | Strong signal (e.g., `data-reactroot` for React) |
| `0.3–0.5` | Weak/generic signal (e.g., `div#app` for SPA generic) |

Multiple matching signals increase the score via a hit multiplier. Confidence is computed as `1 − e^(−score)` and capped at `1.0`. When a meta-framework (e.g., Next.js) is detected, implied frameworks (React) are suppressed unless they have independent high-confidence signals.

### Categories

| Category | Description |
|---|---|
| `spa_framework` | Frontend libraries: React, Vue, Angular, Svelte, etc. |
| `meta_framework` | Full-stack frameworks: Next.js, Nuxt, Remix, Astro, etc. |
| `hosting` | Platform detection: Vercel, Netlify, Cloudflare Pages, etc. |
| `bundler` | Build tools: Webpack, Vite, esbuild, etc. |

### Matcher Types

| Matcher | Description |
|---|---|
| `HeaderContains` | Response header contains substring |
| `HeaderRegex` | Response header matches regex |
| `AnyHeader` | Response header is present |
| `CookieNameHas` | Cookie with given name exists |
| `HTMLSubstring` | Body contains substring (case-insensitive) |
| `HTMLRegex` | Body matches regex |
| `Selector` | CSS selector matches elements |
| `SelectorAttr` | Element matching selector has attribute containing substring |
| `AttrExists` | Any element has the given attribute |
| `ScriptSrcContains` | `<script src>` contains substring |
| `LinkHrefContains` | `<link href>` contains substring |
| `MetaName` | `<meta name>` content contains substring |
| `JSGlobal` | JS sandbox found `window.__key` global |
| `Any(ms...)` | At least one matcher fires |
| `All(ms...)` | All matchers must fire |

## Running Tests

```bash
# Unit tests (offline)
go test -v -race -short ./...

# Integration tests (requires network)
go test -v -run TestIntegration -timeout 120s ./...
```

## Example CLI

```bash
go run examples/cmd/main.go https://nextjs.org
```

## License

MIT