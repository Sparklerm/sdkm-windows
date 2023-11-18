package jdk

import (
	"fmt"
	"os"
	"path/filepath"
)

// RemoveJdk 移除指定版本的JDK
func RemoveJdk(version, jdkDir string) {
	fmt.Println("Removing JDK version:", version)
	jdkPath := filepath.Join(jdkDir, version)
	err := os.RemoveAll(jdkPath)
	if err != nil {
		fmt.Println("Error removing JDK:", err)
		os.Exit(1)
	}
	fmt.Println("JDK removed successfully.")
}
