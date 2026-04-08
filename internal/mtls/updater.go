// file: internal/mtls/updater.go
// version: 1.0.0

package mtls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const defaultReleaseURL = "https://api.github.com/repos/jdfalk/mtls-bridge/releases/latest"

// UpdateCheckInfo is persisted to update-check.json.
type UpdateCheckInfo struct {
	LastCheck time.Time `json:"last_check"`
	Version   string    `json:"version"`
}

// UpdateResult is returned by CheckForUpdate.
type UpdateResult struct {
	Available     bool
	LatestVersion string
	AssetURL      string
	ChecksumURL   string
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func parseGitHubRelease(data []byte) (*githubRelease, error) {
	var release githubRelease
	if err := json.Unmarshal(data, &release); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}
	return &release, nil
}

// needsUpdate compares current version against latest.
// Returns false for dev builds (never auto-update non-release builds).
func needsUpdate(current, latest string) bool {
	if current == "dev" || current == "" {
		return false
	}
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")
	if current == latest {
		return false
	}
	cp := strings.Split(current, ".")
	lp := strings.Split(latest, ".")
	for i := 0; i < 3; i++ {
		var c, l int
		if i < len(cp) {
			fmt.Sscanf(cp[i], "%d", &c)
		}
		if i < len(lp) {
			fmt.Sscanf(lp[i], "%d", &l)
		}
		if l > c {
			return true
		}
		if c > l {
			return false
		}
	}
	return false
}

// CheckForUpdate queries the GitHub Releases API for a newer version.
// Pass "" for releaseURL to use the default.
func CheckForUpdate(currentVersion, releaseURL string) (*UpdateResult, error) {
	if releaseURL == "" {
		releaseURL = defaultReleaseURL
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(releaseURL)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	release, err := parseGitHubRelease(body)
	if err != nil {
		return nil, err
	}

	result := &UpdateResult{
		Available:     needsUpdate(currentVersion, release.TagName),
		LatestVersion: release.TagName,
	}

	if result.Available {
		assetName := assetNameForPlatform()
		for _, asset := range release.Assets {
			if asset.Name == assetName {
				result.AssetURL = asset.BrowserDownloadURL
			}
			if asset.Name == "checksums.txt" {
				result.ChecksumURL = asset.BrowserDownloadURL
			}
		}
	}

	return result, nil
}

func assetNameForPlatform() string {
	osName := strings.ToUpper(runtime.GOOS[:1]) + runtime.GOOS[1:]
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("mtls-bridge_%s_%s.%s", osName, arch, ext)
}

// SelfUpdate downloads and replaces the current binary.
func SelfUpdate(assetURL string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(assetURL)
	if err != nil {
		return fmt.Errorf("download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	tmpFile := execPath + ".update"
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpFile)
		return fmt.Errorf("write update: %w", err)
	}
	f.Close()

	if err := os.Rename(tmpFile, execPath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("replace binary: %w", err)
	}

	return nil
}

// ShouldSkipUpdateCheck returns true if the last check was within the throttle duration.
func (d *Dir) ShouldSkipUpdateCheck(throttle time.Duration) bool {
	data, err := os.ReadFile(d.Path("update-check.json"))
	if err != nil {
		return false
	}
	var info UpdateCheckInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return false
	}
	return time.Since(info.LastCheck) < throttle
}

// WriteUpdateCheck records that an update check was performed.
func (d *Dir) WriteUpdateCheck(version string) error {
	if err := d.EnsureDir(); err != nil {
		return err
	}
	info := UpdateCheckInfo{
		LastCheck: time.Now(),
		Version:   version,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return os.WriteFile(d.Path("update-check.json"), data, 0644)
}
