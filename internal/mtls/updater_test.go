// file: internal/mtls/updater_test.go
// version: 1.0.0

package mtls

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitHubRelease(t *testing.T) {
	body := `{
		"tag_name": "v1.2.3",
		"assets": [
			{"name": "mtls-bridge_1.2.3_Darwin_arm64.tar.gz", "browser_download_url": "https://example.com/darwin-arm64.tar.gz"},
			{"name": "checksums.txt", "browser_download_url": "https://example.com/checksums.txt"}
		]
	}`
	release, err := parseGitHubRelease([]byte(body))
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3", release.TagName)
	assert.Len(t, release.Assets, 2)
}

func TestNeedsUpdate(t *testing.T) {
	assert.True(t, needsUpdate("v1.0.0", "v1.1.0"))
	assert.True(t, needsUpdate("v1.0.0", "v2.0.0"))
	assert.False(t, needsUpdate("v1.1.0", "v1.1.0"))
	assert.False(t, needsUpdate("v1.2.0", "v1.1.0"))
	assert.False(t, needsUpdate("dev", "v1.0.0"))
}

func TestCheckForUpdate_NewVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tag_name": "v2.0.0",
			"assets":   []interface{}{},
		})
	}))
	defer server.Close()

	result, err := CheckForUpdate("v1.0.0", server.URL)
	require.NoError(t, err)
	assert.True(t, result.Available)
	assert.Equal(t, "v2.0.0", result.LatestVersion)
}

func TestCheckForUpdate_AlreadyCurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tag_name": "v1.0.0",
			"assets":   []interface{}{},
		})
	}))
	defer server.Close()

	result, err := CheckForUpdate("v1.0.0", server.URL)
	require.NoError(t, err)
	assert.False(t, result.Available)
}

func TestUpdateCheck_Throttle(t *testing.T) {
	dir := t.TempDir()
	d := NewDir(dir)

	check := UpdateCheckInfo{LastCheck: time.Now(), Version: "v1.0.0"}
	data, _ := json.Marshal(check)
	os.WriteFile(filepath.Join(dir, "update-check.json"), data, 0644)
	assert.True(t, d.ShouldSkipUpdateCheck(1*time.Hour))

	check.LastCheck = time.Now().Add(-2 * time.Hour)
	data, _ = json.Marshal(check)
	os.WriteFile(filepath.Join(dir, "update-check.json"), data, 0644)
	assert.False(t, d.ShouldSkipUpdateCheck(1*time.Hour))
}
