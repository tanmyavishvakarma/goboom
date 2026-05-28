package server

import (
	"fmt"
	"net/http"

	"goboom/internal"

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
	id := internal.GenerateUniqueID()
	directory := fmt.Sprintf("tmp/output%s", id)
	if err := internal.CloneRepoInDirectory(repoUrl, directory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clone repository"})
		return
	}

	c.JSON(http.StatusOK, CloneRepoResponse{ID: id, Directory: directory})
}
