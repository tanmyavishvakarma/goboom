package internal

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func GenerateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func CloneRepoInDirectory(repoUrl string, directory string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	cloneCmd := exec.Command("git", "clone", repoUrl, ".")
	cloneCmd.Dir = directory
	return cloneCmd.Run()
}
