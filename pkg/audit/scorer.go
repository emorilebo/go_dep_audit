package audit

import (
	"math"
	"time"
)

// CalculateHealthScore computes a 0-100 health score for a module
func CalculateHealthScore(metadata *ModuleMetadata, config ScoringConfig) int {
	if metadata == nil {
		return 0
	}

	recencyScore := calculateRecencyScore(metadata.LastCommitDate)
	// For version frequency, we'd need more historical data, but let's use VersionCount as a proxy for maturity
	// or if we had commit frequency.
	// Let's assume we have some basic metrics.
	
	// If we don't have repo metadata, we rely heavily on recency and version count
	versionScore := calculateVersionScore(metadata.VersionCount)
	
	// If we have repo metadata
	commitScore := 0.0
	communityScore := 0.0
	
	if metadata.RepositoryURL != "" {
		commitScore = calculateCommitScore(metadata.CommitFrequency)
		communityScore = calculateCommunityScore(metadata.Stars, metadata.Contributors)
	} else {
		// If no repo metadata, re-distribute weights or just use what we have?
		// For now, let's just average what we have, effectively treating missing as 0 or neutral?
		// Better to treat missing as neutral (50) or ignore the weight?
		// Let's treat missing as 0 for "Community" (no signal) but maybe neutral for activity if we only have publish date.
		// Actually, if we only have proxy info, we only have Recency and VersionCount.
		commitScore = 50 // Neutral assumption
		communityScore = 0 // No signal
	}

	totalScore := (float64(recencyScore) * config.RecencyWeight) +
		(float64(versionScore) * config.VersionFreqWeight) +
		(commitScore * config.CommitActivityWeight) +
		(communityScore * config.CommunityWeight)

	// Normalize to 0-100 int
	score := int(math.Round(totalScore))
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

func calculateRecencyScore(lastDate time.Time) int {
	if lastDate.IsZero() {
		return 0
	}
	hoursSince := time.Since(lastDate).Hours()
	daysSince := hoursSince / 24.0

	// Exponential decay: 
	// 0 days = 100
	// 180 days (6 months) = ~36
	// 365 days (1 year) = ~13
	score := 100 * math.Exp(-daysSince/180.0)
	return int(score)
}

func calculateVersionScore(count int) int {
	// More versions generally means more maturity, up to a point
	// 0 = 0
	// 10+ = 100
	if count >= 20 {
		return 100
	}
	return count * 5
}

func calculateCommitScore(commitsPerMonth float64) float64 {
	// 0 = 0
	// 10/month = 100
	score := commitsPerMonth * 10
	if score > 100 {
		return 100
	}
	return score
}

func calculateCommunityScore(stars, contributors int) float64 {
	// Simple heuristic
	// 1000 stars = 100
	// 50 contributors = 100
	
	starScore := float64(stars) / 10.0 // 1000 stars -> 100
	contribScore := float64(contributors) * 2.0 // 50 contribs -> 100
	
	score := (starScore + contribScore) / 2.0
	if score > 100 {
		return 100
	}
	return score
}

// CategorizeHealth maps a score to a health category
func CategorizeHealth(score int, config ScoringConfig) HealthCategory {
	if score >= config.HealthyThreshold {
		return Healthy
	}
	if score >= config.WarningThreshold {
		return Warning
	}
	if score >= config.StaleThreshold {
		return Stale
	}
	return Risky
}
