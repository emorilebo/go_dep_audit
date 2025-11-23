package audit

import "time"

// AuditConfig holds the configuration for the audit process
type AuditConfig struct {
	ProjectPath       string        `json:"project_path" yaml:"project_path"`
	IncludeIndirect   bool          `json:"include_indirect" yaml:"include_indirect"`
	FetchRepoMetadata bool          `json:"fetch_repo_metadata" yaml:"fetch_repo_metadata"`
	CacheDir          string        `json:"cache_dir" yaml:"cache_dir"`
	CacheTTL          time.Duration `json:"cache_ttl" yaml:"cache_ttl"`

	// Scoring weights and thresholds
	Scoring ScoringConfig `json:"scoring" yaml:"scoring"`

	// License policy
	LicensePolicy LicensePolicy `json:"license_policy" yaml:"license_policy"`

	// Ignore patterns
	IgnoreModules []string `json:"ignore_modules" yaml:"ignore_modules"`

	// API tokens
	GitHubToken string `json:"github_token" yaml:"github_token"`
	GitLabToken string `json:"gitlab_token" yaml:"gitlab_token"`
}

// ScoringConfig defines weights and thresholds for health scoring
type ScoringConfig struct {
	RecencyWeight        float64 `json:"recency_weight" yaml:"recency_weight"`
	VersionFreqWeight    float64 `json:"version_freq_weight" yaml:"version_freq_weight"`
	CommitActivityWeight float64 `json:"commit_activity_weight" yaml:"commit_activity_weight"`
	CommunityWeight      float64 `json:"community_weight" yaml:"community_weight"`

	// Thresholds for categories
	HealthyThreshold int `json:"healthy_threshold" yaml:"healthy_threshold"`
	WarningThreshold int `json:"warning_threshold" yaml:"warning_threshold"`
	StaleThreshold   int `json:"stale_threshold" yaml:"stale_threshold"`
}

// DefaultScoringConfig returns the default scoring configuration
func DefaultScoringConfig() ScoringConfig {
	return ScoringConfig{
		RecencyWeight:        0.4,
		VersionFreqWeight:    0.2,
		CommitActivityWeight: 0.2,
		CommunityWeight:      0.2,
		HealthyThreshold:     70,
		WarningThreshold:     50,
		StaleThreshold:       30,
	}
}

// LicensePolicy defines which licenses are allowed or blocked
type LicensePolicy struct {
	AllowedLicenses []string `json:"allowed_licenses" yaml:"allowed_licenses"`
	BlockedLicenses []string `json:"blocked_licenses" yaml:"blocked_licenses"`
	WarnOnCopyleft  bool     `json:"warn_on_copyleft" yaml:"warn_on_copyleft"`
	WarnOnUnknown   bool     `json:"warn_on_unknown" yaml:"warn_on_unknown"`
}

// DefaultLicensePolicy returns a safe default license policy
func DefaultLicensePolicy() LicensePolicy {
	return LicensePolicy{
		AllowedLicenses: []string{"MIT", "Apache-2.0", "BSD-3-Clause", "BSD-2-Clause", "ISC"},
		WarnOnCopyleft:  true,
		WarnOnUnknown:   true,
	}
}
