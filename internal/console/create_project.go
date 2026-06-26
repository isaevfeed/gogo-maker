package console

import (
	"fmt"
	"gogo-maker/internal/config"
	"gogo-maker/internal/models"
	"os"
	"path/filepath"
)

type CreateProjectParams struct {
	ProjectType models.ProjectType
	Name        string
	DestDir     string
}

func (c *Console) CreateProject(params CreateProjectParams) error {
	// Find repo config by type
	var selectedRepo *config.Repo
	for i := range c.cfg.Reps {
		if models.ProjectType(c.cfg.Reps[i].Type) == params.ProjectType {
			selectedRepo = &c.cfg.Reps[i]
			break
		}
	}

	if selectedRepo == nil {
		return fmt.Errorf("repository with type '%s' not found in config", params.ProjectType)
	}

	// Determine destination directory
	destDir := params.DestDir
	if destDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os.Getwd: %w", err)
		}
		destDir = cwd
	}

	// Clone repo to temp directory first
	tempDir, err := os.MkdirTemp("", "go-helper-clone-*")
	if err != nil {
		return fmt.Errorf("os.MkdirTemp: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp dir

	// Clone
	if err := c.gitRepo.Clone(selectedRepo.Url, tempDir); err != nil {
		return fmt.Errorf("gitRepo.Clone: %w", err)
	}

	// Remove .git directory
	if err := c.gitRepo.RemoveGitDir(tempDir); err != nil {
		return fmt.Errorf("gitRepo.RemoveGitDir: %w", err)
	}

	// Extract template name from URL
	templateName := extractRepoName(selectedRepo.Url)

	// Rename module in go.mod and cmd directory
	if err := c.gitRepo.RenameModule(tempDir, templateName, params.Name); err != nil {
		return fmt.Errorf("gitRepo.RenameModule: %w", err)
	}

	// Copy to final destination
	finalDest := filepath.Join(destDir, params.Name)

	// Create destination directory if needed
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll dest: %w", err)
	}

	if err := c.gitRepo.CopyDir(tempDir, finalDest); err != nil {
		return fmt.Errorf("gitRepo.CopyDir: %w", err)
	}

	c.log.Successf("Project '%s' created at %s", params.Name, finalDest)

	return nil
}

func (c *Console) ListAvailableProjects() {
	c.log.Info("Available project types:")
	for _, repo := range c.cfg.Reps {
		c.log.Infof("  %s  %s", repo.Type, repo.Name)
	}
}

func extractRepoName(url string) string {
	// Extract repo name from URL
	// git@github.com:user/repo.git -> repo
	// https://github.com/user/repo.git -> repo
	base := filepath.Base(url)
	// Remove .git suffix
	if len(base) > 4 && base[len(base)-4:] == ".git" {
		return base[:len(base)-4]
	}
	return base
}
