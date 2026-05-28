package server

import (
	"net/http"
	"os"
	"path/filepath"

	"goboom/internal/helper"
	"goboom/internal/storage"

	"github.com/gin-gonic/gin"
)

type CloneRepoRequest struct {
	RepoURL string `json:"repoUrl"`
}

type CloneRepoResponse struct {
	ID        string `json:"id" binding:"required"`
	Directory string `json:"directory,omitempty"`
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) cloneRepo(c *gin.Context) {
	var req CloneRepoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	repoUrl := req.RepoURL
	id := helper.GenerateUniqueID()
	dirName := helper.ExtractRepoName(repoUrl) + "-" + id

	cwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	directory := filepath.Join(cwd, "tmp", dirName)
	if err := helper.CloneRepoInDirectory(repoUrl, directory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	objectKey := filepath.Join("repos", dirName)
	if err := storage.UploadDirectoryToS3(os.Getenv("AWS_S3_BUCKET_NAME"), objectKey, directory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer helper.DeleteDirectory(directory)

	c.JSON(http.StatusOK, CloneRepoResponse{ID: id, Directory: objectKey})
}
