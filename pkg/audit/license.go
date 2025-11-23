package audit

import (
	"context"
	"strings"
)

// DetectLicense attempts to find and classify the module's license
func DetectLicense(ctx context.Context, modulePath, version string) (string, error) {
	// In a real implementation, this would:
	// 1. Check go.mod/pkg.go.dev metadata if available
	// 2. Fetch LICENSE file from the repo or proxy
	// 3. Use a license classifier (e.g., go-license-detector)
	
	// For this MVP, we'll implement a stub that returns "Unknown" or mocks based on path
	// or fetches from proxy zip (too complex for MVP without external libs).
	
	// Let's try to infer from well-known paths or return Unknown
	if strings.Contains(modulePath, "github.com/") {
		// We could fetch the LICENSE file content via raw.githubusercontent.com if we had the repo URL
		// For now, return "Unknown" to be safe, or "MIT" as a placeholder for testing
		return "Unknown", nil
	}
	return "Unknown", nil
}

// ClassifyLicense categorizes a license string into risk levels
func ClassifyLicense(license string) LicenseRisk {
	l := strings.ToLower(license)
	
	switch {
	case strings.Contains(l, "mit"), strings.Contains(l, "apache"), strings.Contains(l, "bsd"), strings.Contains(l, "isc"):
		return LicensePermissive
	case strings.Contains(l, "gpl"), strings.Contains(l, "agpl"), strings.Contains(l, "mozilla"):
		return LicenseCopyleft
	case strings.Contains(l, "proprietary"), strings.Contains(l, "commercial"):
		return LicenseRestrictive
	default:
		return LicenseUnknown
	}
}

// CheckLicensePolicy validates license against policy
func CheckLicensePolicy(license string, policy LicensePolicy) (allowed bool, warnings []string) {
	// Normalize
	l := strings.ToUpper(license)
	
	// Check blocked
	for _, blocked := range policy.BlockedLicenses {
		if strings.Contains(l, strings.ToUpper(blocked)) {
			return false, []string{"License is explicitly blocked"}
		}
	}
	
	// Check allowed (if whitelist mode)
	if len(policy.AllowedLicenses) > 0 {
		found := false
		for _, allowed := range policy.AllowedLicenses {
			if strings.Contains(l, strings.ToUpper(allowed)) {
				found = true
				break
			}
		}
		if !found {
			// If not in allowed list, is it a hard fail? 
			// Usually whitelist implies everything else is blocked, but let's be soft if not explicitly blocked
			// unless we want strict whitelist.
			// Let's assume strict whitelist if AllowedLicenses is populated.
			return false, []string{"License is not in the allowed list"}
		}
	}
	
	risk := ClassifyLicense(license)
	if policy.WarnOnCopyleft && risk == LicenseCopyleft {
		warnings = append(warnings, "Copyleft license detected")
	}
	if policy.WarnOnUnknown && risk == LicenseUnknown {
		warnings = append(warnings, "Unknown license detected")
	}
	
	return true, warnings
}
