package audit

import (
	"testing"
	"time"
)

func TestCalculateHealthScore(t *testing.T) {
	config := DefaultScoringConfig()

	tests := []struct {
		name     string
		metadata *ModuleMetadata
		wantMin  int
		wantMax  int
	}{
		{
			name: "New and Active",
			metadata: &ModuleMetadata{
				LastCommitDate:  time.Now(),
				VersionCount:    25,
				CommitFrequency: 10,
				Stars:           1000,
				Contributors:    50,
				RepositoryURL:   "https://github.com/example/repo",
			},
			wantMin: 90,
			wantMax: 100,
		},
		{
			name: "Old and Stale",
			metadata: &ModuleMetadata{
				LastCommitDate:  time.Now().AddDate(-2, 0, 0), // 2 years ago
				VersionCount:    5,
				CommitFrequency: 0,
				Stars:           10,
				Contributors:    1,
				RepositoryURL:   "https://github.com/example/repo",
			},
			wantMin: 0,
			wantMax: 40,
		},
		{
			name:     "No Metadata",
			metadata: nil,
			wantMin:  0,
			wantMax:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateHealthScore(tt.metadata, config)
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("CalculateHealthScore() = %v, want between %v and %v", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestCategorizeHealth(t *testing.T) {
	config := DefaultScoringConfig()

	tests := []struct {
		score int
		want  HealthCategory
	}{
		{80, Healthy},
		{60, Warning},
		{40, Stale},
		{20, Risky},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := CategorizeHealth(tt.score, config); got != tt.want {
				t.Errorf("CategorizeHealth(%v) = %v, want %v", tt.score, got, tt.want)
			}
		})
	}
}
