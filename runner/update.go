package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type updateCache struct {
	LastChecked   time.Time `json:"last_checked"`
	LatestVersion string    `json:"latest_version"`
}

func getCacheFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, "gradle-for-agents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "update_cache.json"), nil
}

func loadUpdateCache() updateCache {
	path, err := getCacheFilePath()
	if err != nil {
		return updateCache{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return updateCache{}
	}
	var cache updateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return updateCache{}
	}
	return cache
}

func saveUpdateCache(cache updateCache) {
	path, err := getCacheFilePath()
	if err != nil {
		return
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}

func isHomebrew() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}
	resolved, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		resolved = execPath
	}
	return strings.Contains(resolved, "/Cellar/") || strings.Contains(resolved, "/homebrew/") || strings.Contains(resolved, "/linuxbrew/")
}

func isNewerVersion(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	if current == "dev" || current == "" {
		return false
	}

	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		currNum, err1 := strconv.Atoi(currentParts[i])
		latNum, err2 := strconv.Atoi(latestParts[i])
		if err1 == nil && err2 == nil {
			if latNum > currNum {
				return true
			}
			if currNum > latNum {
				return false
			}
		} else {
			if latestParts[i] > currentParts[i] {
				return true
			}
			if currentParts[i] > latestParts[i] {
				return false
			}
		}
	}
	return len(latestParts) > len(currentParts)
}

func fetchLatestVersionAndCache() {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/silverAndroid/gradle-for-agents/releases/latest", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "gradle-for-agents-update-checker")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	saveUpdateCache(updateCache{
		LastChecked:   time.Now(),
		LatestVersion: release.TagName,
	})
}

// CheckForUpdates inspects if an update warning needs to be displayed.
// It skips checks if installed via Homebrew.
func CheckForUpdates(currentVersion string) {
	if isHomebrew() {
		return
	}

	cache := loadUpdateCache()

	if cache.LatestVersion != "" && isNewerVersion(currentVersion, cache.LatestVersion) {
		fmt.Fprintf(os.Stderr, "\nWARNING: A new version of gradle-for-agents is available: %s (current: %s)\n", cache.LatestVersion, currentVersion)
		fmt.Fprintf(os.Stderr, "To update, run: curl -fsSL https://raw.githubusercontent.com/silverAndroid/gradle-for-agents/main/install.sh | bash\n\n")
	}

	if time.Since(cache.LastChecked) > 24*time.Hour {
		go fetchLatestVersionAndCache()
	}
}
