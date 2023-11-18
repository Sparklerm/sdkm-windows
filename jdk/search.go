package jdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	releasesURL = "https://api.github.com/repos/%s/%s/releases"
)

// Release represents a GitHub release
type Asset struct {
	AssetName   string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// SearchGithubVersion 搜索 GitHub 上的 JDK 版本
func SearchGithubVersion(jdkType string) []Release {
	var releases []Release
	if jdkType == "graalvm" {
		releases = searchGraalVMRelease()
		fmt.Println("GraalVM GitHub Releases:")
		for _, release := range releases {
			fmt.Println("gh-graalvm-" + release.TagName)
		}
	} else {
		fmt.Println("Please specify a valid JDK type: graalvm , zulu, adoptopenjdk, amazoncorretto, bellsoft, sapmachine, liberica, temurin")
	}
	return releases
}

// GetGithubReleaseDownloadUrl 获取 GitHub 上的 JDK 下载地址
func GetGithubReleaseDownloadUrl(jdkType, version string) string {
	githubReleaseDownloadUrl := ""
	githubReleases := SearchGithubVersion(jdkType)
	for _, release := range githubReleases {
		if strings.HasPrefix(release.TagName, version) {
			for _, asset := range release.Assets {
				if strings.Contains(asset.AssetName, "windows") && strings.HasSuffix(asset.AssetName, ".zip") {
					githubReleaseDownloadUrl = asset.DownloadUrl
				}
			}
		}
	}
	return githubReleaseDownloadUrl
}

// searchGraalVMRelease  搜索 GraalVM 的 GitHub Release
func searchGraalVMRelease() []Release {
	// 替换为 GitHub 仓库的所有者（用户名或组织名）
	graalvmUrl := fmt.Sprintf(releasesURL, "graalvm", "graalvm-ce-builds")
	releases, err := getGitHubReleases(graalvmUrl)
	if err != nil {
		fmt.Printf("Error getting GitHub releases: %v\n", err)
		return nil
	}

	var result []Release

	for _, release := range releases {
		if strings.HasPrefix(release.TagName, "jdk") {
			result = append(result, release)
		}
	}
	return result
}

// getGitHubReleases 获取 GitHub Release
func getGitHubReleases(repository string) ([]Release, error) {
	resp, err := http.Get(repository)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}
