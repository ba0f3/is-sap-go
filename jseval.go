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
	_ = vm.Set("window", vm.GlobalObject())

	// document stubs
	docObj := vm.NewObject()
	_ = docObj.Set("createElement", func(tag string) goja.Value { return vm.NewObject() })
	_ = docObj.Set("getElementById", func(id string) goja.Value { return nil })
	_ = docObj.Set("querySelector", func(sel string) goja.Value { return nil })
	_ = docObj.Set("querySelectorAll", func(sel string) goja.Value { return vm.NewArray() })
	_ = docObj.Set("addEventListener", func(...interface{}) {})
	_ = vm.Set("document", docObj)

	// navigator stubs
	navObj := vm.NewObject()
	_ = navObj.Set("userAgent", "Mozilla/5.0")
	_ = vm.Set("navigator", navObj)

	// localStorage / sessionStorage
	storageObj := vm.NewObject()
	_ = storageObj.Set("getItem", func(key string) goja.Value { return nil })
	_ = storageObj.Set("setItem", func(...interface{}) {})
	_ = vm.Set("localStorage", storageObj)
	_ = vm.Set("sessionStorage", storageObj)

	// console stubs
	consoleObj := vm.NewObject()
	_ = consoleObj.Set("log", func(...interface{}) {})
	_ = consoleObj.Set("error", func(...interface{}) {})
	_ = vm.Set("console", consoleObj)

	// fetch stub
	_ = vm.Set("fetch", func(...interface{}) goja.Value {
		return vm.NewObject()
	})

	// setTimeout/setInterval stubs
	_ = vm.Set("setTimeout", func(...interface{}) {})
	_ = vm.Set("setInterval", func(...interface{}) {})
	_ = vm.Set("clearTimeout", func(...interface{}) {})
	_ = vm.Set("clearInterval", func(...interface{}) {})

	// requestAnimationFrame
	_ = vm.Set("requestAnimationFrame", func(...interface{}) {})
	_ = vm.Set("cancelAnimationFrame", func(...interface{}) {})

	// MutationObserver
	_ = vm.Set("MutationObserver", func(...interface{}) goja.Value {
		return vm.NewObject()
	})

	// IntersectionObserver
	_ = vm.Set("IntersectionObserver", func(...interface{}) goja.Value {
		return vm.NewObject()
	})

	// ResizeObserver
	_ = vm.Set("ResizeObserver", func(...interface{}) goja.Value {
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
