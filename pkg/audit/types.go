package audit

import "time"

// HealthCategory represents the overall health status of a dependency
type HealthCategory int

const (
	Healthy HealthCategory = iota
	Warning
	Stale
	Risky
)

func (h HealthCategory) String() string {
	switch h {
	case Healthy:
		return "Healthy"
	case Warning:
		return "Warning"
	case Stale:
		return "Stale"
	case Risky:
		return "Risky"
	default:
		return "Unknown"
	}
}

// LicenseRisk represents the risk level associated with a license
type LicenseRisk int

const (
	LicensePermissive LicenseRisk = iota
	LicenseCopyleft
	LicenseRestrictive
	LicenseUnknown
)

func (l LicenseRisk) String() string {
	switch l {
	case LicensePermissive:
		return "Permissive"
	case LicenseCopyleft:
		return "Copyleft"
	case LicenseRestrictive:
		return "Restrictive"
	case LicenseUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// ModuleHealth contains the audit results for a single module
type ModuleHealth struct {
	Path           string          `json:"path"`
	Version        string          `json:"version"`
	HealthScore    int             `json:"health_score"` // 0-100
	HealthCategory HealthCategory  `json:"health_category"`
	License        string          `json:"license"`
	LicenseRisk    LicenseRisk     `json:"license_risk"`
	FootprintRisk  float64         `json:"footprint_risk"`
	LastPublished  time.Time       `json:"last_published"`
	TransitiveDeps int             `json:"transitive_deps"`
	DirectDep      bool            `json:"direct_dep"`
	Metadata       *ModuleMetadata `json:"metadata,omitempty"`
}

// ModuleMetadata contains raw metadata fetched from sources
type ModuleMetadata struct {
	RepositoryURL   string    `json:"repository_url"`
	Stars           int       `json:"stars"`
	Forks           int       `json:"forks"`
	OpenIssues      int       `json:"open_issues"`
	LastCommitDate  time.Time `json:"last_commit_date"`
	CommitFrequency float64   `json:"commit_frequency"` // commits per month
	Contributors    int       `json:"contributors"`
	VersionCount    int       `json:"version_count"`
}
