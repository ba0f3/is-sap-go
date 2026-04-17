// +build ignore

// This is an example CLI for testing is-sap-go detection.
// Build with: go run examples/cmd/main.go <url>

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ba0f3/is-sap-go"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Request timeout")
	jsEval := flag.Bool("js-eval", false, "Enable JS evaluation (goja sandbox)")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: go run examples/cmd/main.go [options] <url>\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	result, err := sap.Detect(ctx, url,
		sap.WithTimeout(*timeout),
		sap.WithEnableJSEval(*jsEval),
		sap.WithUserAgent("is-sap-go-example/1.0"),
	)
	if err != nil {
		log.Fatalf("Detection failed: %v", err)
	}

	// Marshal to JSON and print
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %v", err)
	}

	fmt.Println(string(data))
}
