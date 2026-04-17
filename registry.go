package sap

import (
	"math"
	"sort"
)

// getAllSignatures returns all registered signatures in priority order.
func getAllSignatures() []*Signature {
	var all []*Signature
	all = append(all, metaFrameworkSignatures...)
	all = append(all, frontendSignatures...)
	all = append(all, hostingSignatures...)
	all = append(all, genericSignatures...)
	return all
}

// scoreSignatures evaluates all registered signatures against the scan and populates the Result.
func scoreSignatures(scan *Scan, result *Result) {
	sigs := getAllSignatures()

	// Score each signature.
	type SigScore struct {
		sig       *Signature
		hits      int
		evidences []string
		score     float64
		confidence float64
	}

	var scored []*SigScore
	for _, sig := range sigs {
		var hits int
		var evidences []string

		// Evaluate each matcher in the signature; any match counts.
		for _, m := range sig.Matchers {
			if hit, ev := m(scan); hit {
				hits++
				if ev != "" {
					evidences = append(evidences, ev)
				}
			}
		}

		if hits == 0 {
			continue // No match for this signature.
		}

		// Calculate score: weight * (1 + 0.25 * (hits - 1)), capped at weight * 2
		score := sig.Weight * (1 + 0.25*float64(hits-1))
		if score > sig.Weight*2 {
			score = sig.Weight * 2
		}

		// Calculate confidence: 1 - exp(-score), capped at 1.0
		confidence := 1 - math.Exp(-score)
		if confidence > 1.0 {
			confidence = 1.0
		}

		scored = append(scored, &SigScore{
			sig:        sig,
			hits:       hits,
			evidences:  evidences,
			score:      score,
			confidence: confidence,
		})
	}

	// Sort by confidence descending.
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].confidence > scored[j].confidence
	})

	// Populate Result with detected frameworks.
	for _, ss := range scored {
		fw := Framework{
			Name:       ss.sig.Framework,
			Category:   ss.sig.Category,
			Confidence: ss.confidence,
			Signals:    ss.evidences,
		}

		// Extract version if available.
		if ss.sig.Version != nil {
			fw.Version = ss.sig.Version(scan)
		}

		// Categorize: hosting or framework/meta.
		if ss.sig.Category == CategoryHosting {
			result.Hosting = append(result.Hosting, fw)
		} else {
			result.Frameworks = append(result.Frameworks, fw)
		}
	}
}

// suppressRedundant removes redundant frontend frameworks when a meta-framework is high-confidence.
// E.g., suppress React if Next.js is detected with confidence ≥ 0.6.
func suppressRedundant(result *Result) {
	type ImpliesMap map[string]bool

	// Build a map of what's implied by high-confidence frameworks.
	implies := make(ImpliesMap)
	for _, fw := range result.Frameworks {
		if fw.Confidence >= 0.6 {
			sig := lookupSignature(fw.Name)
			if sig != nil {
				for _, imp := range sig.Implies {
					implies[imp] = true
				}
			}
		}
	}

	if len(implies) == 0 {
		return // Nothing to suppress.
	}

	// Filter: keep framework if it's not implied, or has independent strong evidence.
	var filtered []Framework
	for _, fw := range result.Frameworks {
		if !implies[fw.Name] {
			filtered = append(filtered, fw)
			continue
		}

		// Framework is implied; check for independent evidence.
		// For now, we suppress it unless we later add special logic.
		// (In a full implementation, you might check signal count or confidence threshold.)
		// Keep it if confidence is >= 0.8 (strong independent evidence).
		if fw.Confidence >= 0.8 {
			filtered = append(filtered, fw)
		}
	}

	result.Frameworks = filtered
}

// sortFrameworks sorts frameworks by confidence descending.
func sortFrameworks(fws []Framework) {
	sort.Slice(fws, func(i, j int) bool {
		if fws[i].Confidence != fws[j].Confidence {
			return fws[i].Confidence > fws[j].Confidence
		}
		return fws[i].Name < fws[j].Name
	})
}

// determineSPA checks if any framework indicates a SPA.
func determineSPA(result *Result) bool {
	// If any CategorySPAFramework or CategoryMetaFramework has confidence >= 0.3, it's a SPA.
	// (Lower threshold to catch single strong signals)
	for _, fw := range result.Frameworks {
		if (fw.Category == CategorySPAFramework || fw.Category == CategoryMetaFramework) &&
			fw.Confidence >= 0.3 {
			return true
		}
	}

	// TODO: Add generic SPA heuristics here (noscript check, hashed bundles, etc.).
	// For now, rely on explicit framework detection.

	return false
}

// calculateConfidence returns the overall SPA confidence.
func calculateConfidence(result *Result) float64 {
	maxConf := 0.0
	for _, fw := range result.Frameworks {
		if fw.Category == CategorySPAFramework || fw.Category == CategoryMetaFramework {
			if fw.Confidence > maxConf {
				maxConf = fw.Confidence
			}
		}
	}
	return maxConf
}

// lookupSignature finds a Signature by framework name.
func lookupSignature(name string) *Signature {
	for _, sig := range getAllSignatures() {
		if sig.Framework == name {
			return sig
		}
	}
	return nil
}
