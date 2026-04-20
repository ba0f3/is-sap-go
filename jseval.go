package sap

import (
	"regexp"
	"strings"

	"github.com/dop251/goja"
)

// JSBlob represents a harvested window.* global variable.
type JSBlob struct {
	Name string
	Size int
}

// ExtractJSDataFromHTML uses goja to sandbox inline scripts and harvest window.__* globals.
func ExtractJSDataFromHTML(html string) []JSBlob {
	vm := goja.New()

	// Set up minimal sandboxed globals: window, document, navigator, etc.
	// These are stubs to prevent script errors; they won't execute meaningful DOM operations.
	setupSandbox(vm)

	// Extract all inline <script> tags and concatenate their content.
	inlineScripts := extractInlineScripts(html)
	for _, script := range inlineScripts {
		// Try to run each script; ignore errors (scripts may reference undefined globals).
		_, _ = vm.RunString(script)
	}

	// Harvest all window.__* globals.
	var blobs []JSBlob
	globalObj := vm.GlobalObject()
	if globalObj != nil {
		for _, key := range globalObj.Keys() {
			if strings.HasPrefix(key, "__") {
				val := globalObj.Get(key)
				if val != nil {
					sizeEst := estimateSize(val)
					blobs = append(blobs, JSBlob{
						Name: key,
						Size: sizeEst,
					})
				}
			}
		}
	}

	return blobs
}

// setupSandbox initializes minimal stubs to prevent script errors.
func setupSandbox(vm *goja.Runtime) {
	// window object (goja's GlobalObject IS window)
	vm.Set("window", vm.GlobalObject())

	// document stubs
	docObj := vm.NewObject()
	docObj.Set("createElement", func(tag string) goja.Value { return vm.NewObject() })
	docObj.Set("getElementById", func(id string) goja.Value { return nil })
	docObj.Set("querySelector", func(sel string) goja.Value { return nil })
	docObj.Set("querySelectorAll", func(sel string) goja.Value { return vm.NewArray() })
	docObj.Set("addEventListener", func(...interface{}) {})
	vm.Set("document", docObj)

	// navigator stubs
	navObj := vm.NewObject()
	navObj.Set("userAgent", "Mozilla/5.0")
	vm.Set("navigator", navObj)

	// localStorage / sessionStorage
	storageObj := vm.NewObject()
	storageObj.Set("getItem", func(key string) goja.Value { return nil })
	storageObj.Set("setItem", func(...interface{}) {})
	vm.Set("localStorage", storageObj)
	vm.Set("sessionStorage", storageObj)

	// console stubs
	consoleObj := vm.NewObject()
	consoleObj.Set("log", func(...interface{}) {})
	consoleObj.Set("error", func(...interface{}) {})
	vm.Set("console", consoleObj)

	// fetch stub
	vm.Set("fetch", func(...interface{}) goja.Value {
		return vm.NewObject() // return a dummy promise-like object
	})

	// setTimeout/setInterval stubs
	vm.Set("setTimeout", func(...interface{}) {})
	vm.Set("setInterval", func(...interface{}) {})
	vm.Set("clearTimeout", func(...interface{}) {})
	vm.Set("clearInterval", func(...interface{}) {})

	// requestAnimationFrame
	vm.Set("requestAnimationFrame", func(...interface{}) {})
	vm.Set("cancelAnimationFrame", func(...interface{}) {})

	// MutationObserver
	vm.Set("MutationObserver", func(...interface{}) goja.Value {
		return vm.NewObject()
	})

	// IntersectionObserver
	vm.Set("IntersectionObserver", func(...interface{}) goja.Value {
		return vm.NewObject()
	})

	// ResizeObserver
	vm.Set("ResizeObserver", func(...interface{}) goja.Value {
		return vm.NewObject()
	})
}

// extractInlineScripts extracts all inline <script> tag contents from HTML.
func extractInlineScripts(html string) []string {
	// Simple regex to find <script>.....</script> (not type="module" or with src=)
	re := regexp.MustCompile(`<script\s*(?:type\s*=\s*"[^"]*"\s*)?(?:async|defer)?\s*>(.+?)</script>`)

	var scripts []string
	matches := re.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) > 1 {
			scripts = append(scripts, m[1])
		}
	}

	// Also handle data attributes and attribute-based inline scripts more carefully.
	// This is a best-effort extraction; some scripts may be obfuscated or conditional.
	return scripts
}

// estimateSize returns a rough byte size of a goja value.
func estimateSize(val goja.Value) int {
	if val == nil {
		return 0
	}
	s := val.String()
	return len(s)
}
