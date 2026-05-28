package helper

import (
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func GenerateUniqueID() string {
	return uuid.New().String()
}

func CloneRepoInDirectory(repoUrl string, directory string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	cloneCmd := exec.Command("git", "clone", repoUrl, ".")
	cloneCmd.Dir = directory
	return cloneCmd.Run()
}

func ExtractRepoName(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return ""
	}

	repo := strings.TrimPrefix(u.Path, "/")
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.ReplaceAll(repo, "/", "-")
	return repo
}

func DeleteDirectory(directory string) error {
	return os.RemoveAll(directory)
}
