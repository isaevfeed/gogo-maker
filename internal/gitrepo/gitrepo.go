package gitrepo

import (
	"fmt"
	"gogo-maker/internal/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type GitRepo struct {
	log *logger.Logger
}

func New(log *logger.Logger) *GitRepo {
	return &GitRepo{log: log}
}

// Clone clones a repository from URL to destination directory
func (g *GitRepo) Clone(url, dest string) error {
	g.log.Infof("Cloning %s", url)

	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("git.PlainClone: %w", err)
	}

	g.log.Success("Repository cloned")

	return nil
}

// RenameModule renames cmd/{oldName} to cmd/{newName} and updates all imports
func (g *GitRepo) RenameModule(projectPath, oldName, newName string) error {
	// 1. Rename cmd/{oldName} to cmd/{newName} if exists
	cmdOldPath := filepath.Join(projectPath, "cmd", oldName)
	cmdNewPath := filepath.Join(projectPath, "cmd", newName)

	if _, err := os.Stat(cmdOldPath); err == nil {
		g.log.Infof("Renaming cmd/%s -> cmd/%s", oldName, newName)
		if err := os.Rename(cmdOldPath, cmdNewPath); err != nil {
			return fmt.Errorf("os.Rename cmd: %w", err)
		}
	}

	// 2. Update all .go files imports
	if err := g.updateImports(projectPath, oldName, newName); err != nil {
		return fmt.Errorf("updateImports: %w", err)
	}

	g.log.Infof("Updated module: %s -> %s", oldName, newName)

	return nil
}

// shouldProcessFile returns true if the file should be processed for module name replacement
func shouldProcessFile(filename string) bool {
	// Process Go files
	if filepath.Ext(filename) == ".go" {
		return true
	}

	// Process specific files by name
	textFiles := map[string]bool{
		"go.mod":              true,
		"go.sum":              true,
		"Dockerfile":          true,
		"Makefile":            true,
		"docker-compose.yml":  true,
		"docker-compose.yaml": true,
		".dockerignore":       true,
		".gitignore":          true,
		"README.md":           true,
	}

	if textFiles[filename] {
		return true
	}

	// Process files by extension
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml", ".json", ".toml", ".sh", ".md", ".txt":
		return true
	}

	// Handle Dockerfile.* variants (e.g., Dockerfile.prod, Dockerfile.dev)
	if strings.HasPrefix(filename, "Dockerfile") {
		return true
	}

	return false
}

func (g *GitRepo) updateImports(projectPath, oldName, newName string) error {
	return filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip .git and vendor directories
			if info.Name() == ".git" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		filename := filepath.Base(path)
		if !shouldProcessFile(filename) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("os.ReadFile %s: %w", path, err)
		}

		newContent := string(content)

		// Replace module name in imports and requires (handles both with and without subpackages)
		// "oldName" -> "newName"
		// "oldName/subpath" -> "newName/subpath"
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf(`"%s"`, oldName), fmt.Sprintf(`"%s"`, newName))
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf(`"%s/`, oldName), fmt.Sprintf(`"%s/`, newName))

		// Replace module line in go.mod
		if filename == "go.mod" {
			newContent = strings.ReplaceAll(newContent, fmt.Sprintf("module %s\n", oldName), fmt.Sprintf("module %s\n", newName))
		}

		// Replace cmd/oldName paths (for Dockerfile, Makefile, etc.)
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf("cmd/%s", oldName), fmt.Sprintf("cmd/%s", newName))
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf("./cmd/%s", oldName), fmt.Sprintf("./cmd/%s", newName))

		// Replace plain oldName references (for README, configs, etc.)
		// Only if not part of a longer word (e.g., avoid replacing "my-golang-http-server-template-new")
		newContent = strings.ReplaceAll(newContent, oldName, newName)

		if newContent != string(content) {
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("os.WriteFile %s: %w", path, err)
			}
		}

		return nil
	})
}

// RemoveGitDir removes .git directory to clean git history
func (g *GitRepo) RemoveGitDir(projectPath string) error {
	gitDir := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		g.log.Info("Removing .git directory")
		if err := os.RemoveAll(gitDir); err != nil {
			return fmt.Errorf("os.RemoveAll .git: %w", err)
		}
	}

	return nil
}

// CopyDir copies a directory recursively
func (g *GitRepo) CopyDir(src, dst string) error {
	g.log.Infof("Copying to %s", dst)

	err := os.CopyFS(dst, os.DirFS(src))
	if err != nil {
		return fmt.Errorf("os.CopyFS: %w", err)
	}

	g.log.Success("Copy completed")

	return nil
}
