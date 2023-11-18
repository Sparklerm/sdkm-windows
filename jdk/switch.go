package jdk

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sdkm/model"
	"sort"
	"strings"
)

// ListLocalJdkVersions 列出可切换的JDK版本
func ListLocalJdkVersions(config model.Config) {
	fmt.Println("可切换的JDK版本：")
	jdkDir := config.JDKDir
	versions, err := getJDKVersions(jdkDir)
	if err != nil {
		fmt.Println("Error listing JDK versions:", err)
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Println(version)
	}
}

// UseVersion 启用指定版本的JDK
func UseVersion(config model.Config, version string) {
	jdkDir := config.JDKDir
	envName := config.EnvName

	jdkPath := filepath.Join(jdkDir, version)

	// 检查JDK版本是否存在
	if !isJDKVersionAvailable(jdkDir, version) {
		fmt.Println("Invalid JDK version or path does not exist.")
		os.Exit(1)
	}

	// 更新环境变量
	registerPath(envName)
	err := updateEnvironment(envName, jdkPath)
	if err != nil {
		fmt.Println("Error updating environment variable:", err)
		os.Exit(1)
	}

	fmt.Printf("JDK version %s 已启用\n\r", version)
}

// getJDKVersions 获取JDK_DIR目录下所有的JDK版本
func getJDKVersions(jdkDir string) ([]string, error) {
	var versions []string

	entries, err := ioutil.ReadDir(jdkDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		//if entry.IsDir() && strings.HasPrefix(entry.Name(), "jdk") {
		//	versions = append(versions, entry.Name())
		//}
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	// 按版本号排序
	sort.Strings(versions)
	return versions, nil
}

// isJDKVersionAvailable 检查指定版本的JDK是否存在
func isJDKVersionAvailable(jdkDir, version string) bool {
	jdkPath := filepath.Join(jdkDir, version)
	_, err := os.Stat(jdkPath)
	return err == nil
}

// updateEnvironment 更新环境变量
func updateEnvironment(envName, jdkPath string) error {
	// 更新JDK环境变量
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("%s", "%s", [System.EnvironmentVariableTarget]::User)`, envName, jdkPath+"\\bin"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// 在Path环境变量中添加JDK环境变量
func registerPath(envName string) {
	// 获取当前环境变量 "PATH" 的值
	path := exec.Command("powershell", "-Command", `[System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)`)
	output, err := path.Output()
	if err != nil {
		fmt.Println("Error getting PATH environment variable:", err)
		os.Exit(1)
	}
	currentPath := string(output)
	currentPath = strings.TrimSpace(currentPath)
	currentPath = strings.Trim(currentPath, "\n")
	currentPath = strings.Trim(currentPath, "\r")
	// 检查路径是否已经在 "PATH" 中
	jdkEnv := "%" + envName + "%"
	if strings.Contains(currentPath, jdkEnv) {
		return
	}

	// 更新环境变量
	// 如果Path最后一个字符不是分号，则添加分号
	if currentPath[len(currentPath)-1] != ';' {
		currentPath += ";"
	}
	newPath := currentPath + jdkEnv
	pathUpdate := exec.Command("powershell", "-Command", fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("Path", "%s", [System.EnvironmentVariableTarget]::User)`, newPath))
	pathUpdate.Stdout = os.Stdout
	pathUpdate.Stderr = os.Stderr
	pathUpdate.Run()
}
