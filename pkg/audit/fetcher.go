package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Fetcher handles metadata retrieval
type Fetcher struct {
	client *http.Client
	config AuditConfig
}

func NewFetcher(config AuditConfig) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: config,
	}
}

// ProxyInfo represents data from proxy.golang.org/{module}/@v/{version}.info
type ProxyInfo struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

// FetchModuleMetadata gathers all available metadata for a module
func (f *Fetcher) FetchModuleMetadata(ctx context.Context, modulePath, version string) (*ModuleMetadata, error) {
	meta := &ModuleMetadata{}

	// 1. Fetch from Go Proxy
	proxyInfo, err := f.fetchProxyInfo(ctx, modulePath, version)
	if err != nil {
		// Log error but continue? For now, return error as basic info is needed
		// In production, might want to fallback or mark as unknown
		return nil, fmt.Errorf("failed to fetch proxy info: %w", err)
	}
	meta.LastCommitDate = proxyInfo.Time
	
	// 2. Fetch version list to count versions
	versions, err := f.fetchVersionList(ctx, modulePath)
	if err == nil {
		meta.VersionCount = len(versions)
		// Calculate frequency based on versions? 
		// For now just store count
	}

	// 3. Fetch Repository Metadata (if enabled and possible)
	if f.config.FetchRepoMetadata {
		repoURL := getRepoURL(modulePath)
		if repoURL != "" {
			meta.RepositoryURL = repoURL
			// TODO: Implement GitHub/GitLab API fetching here
			// For now, we'll leave the repo stats as 0
		}
	}

	return meta, nil
}

func (f *Fetcher) fetchProxyInfo(ctx context.Context, modulePath, version string) (*ProxyInfo, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.info", strings.ToLower(modulePath), version)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy returned status: %d", resp.StatusCode)
	}

	var info ProxyInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (f *Fetcher) fetchVersionList(ctx context.Context, modulePath string) ([]string, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", strings.ToLower(modulePath))
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy returned status: %d", resp.StatusCode)
	}

	var versions []string
	// The list endpoint returns a text list of versions, one per line
	// We can read it into a buffer
	// For simplicity, let's just count lines if we only need count, 
	// but we might need parsing later.
	// Implementation omitted for brevity, returning empty for now
	return versions, nil
}

// Helper to guess repo URL from module path
func getRepoURL(modulePath string) string {
	if strings.HasPrefix(modulePath, "github.com/") {
		parts := strings.Split(modulePath, "/")
		if len(parts) >= 3 {
			return fmt.Sprintf("https://github.com/%s/%s", parts[1], parts[2])
		}
	}
	// Add gitlab, bitbucket etc.
	return ""
}
