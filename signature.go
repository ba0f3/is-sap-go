package sap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Matcher is a function that tests a scan and returns (matched, evidence).
type Matcher func(s *Scan) (bool, string)

// Signature describes one technology's detection fingerprint.
type Signature struct {
	Framework string
	Category  Category
	// Weight: 0.3 weak / 0.6 strong / 0.95 definitive
	Weight   float64
	Matchers []Matcher          // ANY single match fires the signature
	Version  func(*Scan) string // optional version extractor
	// Implies lists framework names this signature's detection supersedes.
	// E.g. Next.js implies React — if Next.js is detected we suppress React
	// unless React has independent high-confidence signals.
	Implies []string
}

// --- Matcher builders ---

// HeaderContains fires when a response header contains the given substring (case-insensitive).
func HeaderContains(name, substr string) Matcher {
	lname := strings.ToLower(name)
	lsub := strings.ToLower(substr)
	return func(s *Scan) (bool, string) {
		v := s.Headers.Get(lname)
		if v == "" {
			// Headers.Get is canonicalized; try the exact case too
			v = s.Headers.Get(name)
		}
		if strings.Contains(strings.ToLower(v), lsub) && v != "" {
			return true, fmt.Sprintf("header %s: %s", name, v)
		}
		return false, ""
	}
}

// HeaderRegex fires when a response header matches the regex.
func HeaderRegex(name string, re *regexp.Regexp) Matcher {
	return func(s *Scan) (bool, string) {
		v := s.Headers.Get(name)
		if m := re.FindString(v); m != "" {
			return true, fmt.Sprintf("header %s matches %s", name, re.String())
		}
		return false, ""
	}
}

// AnyHeader fires when a response header is present (non-empty).
func AnyHeader(name string) Matcher {
	return func(s *Scan) (bool, string) {
		v := s.Headers.Get(name)
		if v != "" {
			return true, fmt.Sprintf("header %s: %s", name, v)
		}
		return false, ""
	}
}

// CookieNameHas fires when a Set-Cookie header contains a cookie with the given name.
func CookieNameHas(name string) Matcher {
	return func(s *Scan) (bool, string) {
		for _, c := range s.Cookies {
			if strings.EqualFold(c.Name, name) {
				return true, fmt.Sprintf("cookie %s", c.Name)
			}
		}
		return false, ""
	}
}

// HTMLSubstring fires when the response body contains the given substring (case-insensitive).
func HTMLSubstring(sub string) Matcher {
	lower := strings.ToLower(sub)
	return func(s *Scan) (bool, string) {
		if strings.Contains(s.BodyLower, lower) {
			return true, fmt.Sprintf("html contains %q", sub)
		}
		return false, ""
	}
}

// HTMLRegex fires when the response body matches the regex.
func HTMLRegex(re *regexp.Regexp) Matcher {
	return func(s *Scan) (bool, string) {
		if m := re.FindString(s.BodyStr); m != "" {
			return true, fmt.Sprintf("html matches /%s/", re.String())
		}
		return false, ""
	}
}

// Selector fires when a CSS selector matches at least one element in the DOM.
func Selector(sel string) Matcher {
	return func(s *Scan) (bool, string) {
		doc := s.Doc()
		if doc == nil {
			return false, ""
		}
		if doc.Find(sel).Length() > 0 {
			return true, fmt.Sprintf("selector %q", sel)
		}
		return false, ""
	}
}

// SelectorAttr fires when a CSS selector matches an element whose attribute contains substr.
func SelectorAttr(sel, attr, substr string) Matcher {
	lower := strings.ToLower(substr)
	return func(s *Scan) (bool, string) {
		doc := s.Doc()
		if doc == nil {
			return false, ""
		}
		found := false
		doc.Find(sel).Each(func(_ int, node *goquery.Selection) {
			if found {
				return
			}
			v, exists := node.Attr(attr)
			if exists && strings.Contains(strings.ToLower(v), lower) {
				found = true
			}
		})
		if found {
			return true, fmt.Sprintf("selector %q attr[%s] contains %q", sel, attr, substr)
		}
		return false, ""
	}
}

// AttrExists fires when any element has the given attribute (regardless of value).
func AttrExists(attr string) Matcher {
	sel := fmt.Sprintf("[%s]", attr)
	return func(s *Scan) (bool, string) {
		doc := s.Doc()
		if doc == nil {
			// Fast path: substring check
			needle := strings.ToLower(attr + "=")
			if strings.Contains(s.BodyLower, needle) {
				return true, fmt.Sprintf("attr %s exists", attr)
			}
			return false, ""
		}
		if doc.Find(sel).Length() > 0 {
			return true, fmt.Sprintf("attr %s exists", attr)
		}
		return false, ""
	}
}

// ScriptSrcContains fires when any <script src="…"> contains the given substring.
func ScriptSrcContains(substr string) Matcher {
	lower := strings.ToLower(substr)
	return func(s *Scan) (bool, string) {
		// Fast path: raw body check
		if !strings.Contains(s.BodyLower, lower) {
			return false, ""
		}
		doc := s.Doc()
		if doc == nil {
			return true, fmt.Sprintf("script src contains %q (fast)", substr)
		}
		var match string
		doc.Find("script[src]").Each(func(_ int, node *goquery.Selection) {
			if match != "" {
				return
			}
			src, _ := node.Attr("src")
			if strings.Contains(strings.ToLower(src), lower) {
				match = src
			}
		})
		if match != "" {
			return true, fmt.Sprintf("script src %q", match)
		}
		return false, ""
	}
}

// LinkHrefContains fires when any <link href="…"> contains the given substring.
func LinkHrefContains(substr string) Matcher {
	lower := strings.ToLower(substr)
	return func(s *Scan) (bool, string) {
		if !strings.Contains(s.BodyLower, lower) {
			return false, ""
		}
		doc := s.Doc()
		if doc == nil {
			return true, fmt.Sprintf("link href contains %q (fast)", substr)
		}
		var match string
		doc.Find("link[href]").Each(func(_ int, node *goquery.Selection) {
			if match != "" {
				return
			}
			href, _ := node.Attr("href")
			if strings.Contains(strings.ToLower(href), lower) {
				match = href
			}
		})
		if match != "" {
			return true, fmt.Sprintf("link href %q", match)
		}
		return false, ""
	}
}

// MetaName fires when <meta name="…"> has content containing the given substring.
func MetaName(name, contentSubstr string) Matcher {
	lname := strings.ToLower(name)
	lsub := strings.ToLower(contentSubstr)
	return func(s *Scan) (bool, string) {
		doc := s.Doc()
		if doc == nil {
			return false, ""
		}
		found := false
		doc.Find("meta[name]").Each(func(_ int, node *goquery.Selection) {
			if found {
				return
			}
			n, _ := node.Attr("name")
			if strings.ToLower(n) != lname {
				return
			}
			content, _ := node.Attr("content")
			if strings.Contains(strings.ToLower(content), lsub) {
				found = true
			}
		})
		if found {
			return true, fmt.Sprintf("meta[name=%q] content contains %q", name, contentSubstr)
		}
		return false, ""
	}
}

// JSGlobal fires when the JS sandbox found a window.__key global.
// Silently returns false when EnableJSEval is off.
func JSGlobal(key string) Matcher {
	return func(s *Scan) (bool, string) {
		for _, b := range s.JSBlobs() {
			if b.Name == key {
				return true, fmt.Sprintf("js global %s (size=%d)", key, b.Size)
			}
		}
		return false, ""
	}
}

// Any fires when at least one of the given matchers fires.
func Any(ms ...Matcher) Matcher {
	return func(s *Scan) (bool, string) {
		for _, m := range ms {
			if hit, ev := m(s); hit {
				return true, ev
			}
		}
		return false, ""
	}
}

// All fires only when every one of the given matchers fires.
func All(ms ...Matcher) Matcher {
	return func(s *Scan) (bool, string) {
		var evidences []string
		for _, m := range ms {
			hit, ev := m(s)
			if !hit {
				return false, ""
			}
			if ev != "" {
				evidences = append(evidences, ev)
			}
		}
		return true, strings.Join(evidences, "; ")
	}
}
