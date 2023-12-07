package jdk

import (
	"encoding/json"
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sdkm/model"
	"strings"
)

const (
	ZULU_API = "https://api.azul.com/metadata/v1/zulu/packages/?java_version=%s&os=windows&arch=x64&archive_type=zip&java_package_type=jdk&latest=true&distro_version=%s&release_status=ga&certifications=tck&page=1&page_size=1"
)

func Install(jdkType []model.JdkType, version, destination string) {
	fmt.Println("Installing JDK version:", version)
	// 获取JDK安装包下载地址
	fmt.Println("Getting download URL...")
	downloadUrl := getDownloadUrl(jdkType, version)
	if downloadUrl == "" {
		fmt.Println("Error: Unsupported JDK version:", version)
		os.Exit(1)
	}
	fmt.Println("Download URL:", downloadUrl)

	// 下载JDK压缩安装包
	fmt.Println("Downloading JDK...")
	err := downloadPackage(downloadUrl, filepath.Join(destination, "jdk.zip"))
	if err != nil {
		fmt.Println("Error downloading JDK:", err)
		return
	}
	fmt.Println("JDK downloaded successfully.")

	// 解压JDK安装包
	fmt.Println("Extracting JDK...")
	err = unzipFile(filepath.Join(destination, "jdk.zip"), destination)
	if err != nil {
		fmt.Println("Error unzipping JDK package:", err)
		os.Exit(1)
	}
	// 删除JDK安装包
	err = os.Remove(filepath.Join(destination, "jdk.zip"))
	fmt.Println("JDK extracted successfully.")
	// 输出JDK安装路径
	fmt.Println("JDK installed at:", destination)
}

// ListJdkVersions 列出所有可用的 JDK 版本
func ListJdkVersions(jdkTypes []model.JdkType, version string) {
	if strings.Compare(version, "common") == 0 {
		for _, jdkType := range jdkTypes {
			for _, version := range jdkType.Versions {
				fmt.Println(jdkType.JdkType + "-" + version.Version)
			}
		}
		return
	}
	SearchGithubVersion(version)
}

// getZuluReleaseDownloadUrl 获取Zulu下载地址
func getZuluReleaseDownloadUrl(version string) string {
	// 请求ZULU_API 获取返回信息
	apiURL := fmt.Sprintf(ZULU_API, version, version)

	// 创建一个HTTP客户端
	client := &http.Client{}
	// 创建一个GET请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.Status)
		return ""
	}

	// 解码JSON响应
	var responseData []model.ZuluApiResult
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	return responseData[0].DownloadUrl
}

// getDownloadUrl 获取JDK安装包下载地址
func getDownloadUrl(jdkTypes []model.JdkType, version string) string {
	downloadUrl := ""
	versionSplit := strings.Split(version, "-")
	if versionSplit[0] == "gh" {
		// 从 GitHub Releases 获取下载地址
		downloadUrl = GetGithubReleaseDownloadUrl(versionSplit[1], versionSplit[2])
	} else {
		javaVersion := versionSplit[1]
		if strings.HasPrefix(strings.ToLower(version), "oracle") {
			for i := range jdkTypes {
				if strings.ToLower(jdkTypes[i].JdkType) == "oracle" {
					for j := range jdkTypes[i].Versions {
						if jdkTypes[i].Versions[j].Version == javaVersion {
							downloadUrl = jdkTypes[i].Versions[j].DownloadUrl
						}
					}
				}
			}
		} else if strings.HasPrefix(strings.ToLower(version), "graalvm") {
			for i := range jdkTypes {
				if strings.ToLower(jdkTypes[i].JdkType) == "graalvm" {
					for j := range jdkTypes[i].Versions {
						if jdkTypes[i].Versions[j].Version == javaVersion {
							downloadUrl = jdkTypes[i].Versions[j].DownloadUrl
						}
					}
				}
			}
		} else if strings.HasPrefix(strings.ToLower(version), "zulu") {
			zuluReleaseDownloadUrl := getZuluReleaseDownloadUrl(javaVersion)
			downloadUrl = zuluReleaseDownloadUrl
		}
	}

	if downloadUrl == "" || downloadUrl == "null" {
		fmt.Println("Error: Unsupported JDK version:", version)
		return ""
	}

	return downloadUrl
}

// downloadPackage 下载 JDK 安装包
func downloadPackage(downloadUrl, destination string) error {
	response, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

func unzipFile(zipFile, destination string) error {
	return archiver.Unarchive(zipFile, destination)
}
