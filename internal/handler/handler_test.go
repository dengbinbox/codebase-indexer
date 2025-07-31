package handler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"codebase-indexer/internal/repository"
	"codebase-indexer/internal/service"
	"codebase-indexer/test/mocks"
)

func TestNewGRPCHandler(t *testing.T) {
	var mockLogger = &mocks.MockLogger{}
	// Create test objects
	httpSync := &repository.HTTPSync{}
	fileScanner := &repository.FileScanner{}
	storageManager := &repository.StorageManager{}
	scheduler := &service.Scheduler{}

	h := NewGRPCHandler(httpSync, fileScanner, storageManager, scheduler, mockLogger)
	assert.NotNil(t, h)
}

func TestIsGitRepository(t *testing.T) {
	var mockLogger = &mocks.MockLogger{}
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create .git directory
	err := os.Mkdir(filepath.Join(tmpDir, ".git"), 0755)
	assert.NoError(t, err)

	httpSync := &repository.HTTPSync{}
	fileScanner := &repository.FileScanner{}
	storageManager := &repository.StorageManager{}
	scheduler := &service.Scheduler{}
	h := NewGRPCHandler(httpSync, fileScanner, storageManager, scheduler, mockLogger)

	// Test valid git repository
	assert.True(t, h.isGitRepository(tmpDir))

	// Test invalid path
	assert.False(t, h.isGitRepository(filepath.Join(tmpDir, "nonexistent")))

	// Test non-git directory
	nonGitDir := filepath.Join(tmpDir, "not-git")
	err = os.Mkdir(nonGitDir, 0755)
	assert.NoError(t, err)
	assert.False(t, h.isGitRepository(nonGitDir))
}

func TestFindCodebasePathsToRegister(t *testing.T) {
	var mockLogger = &mocks.MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	// Create test directory structure
	baseDir := t.TempDir()

	// Create subdirectory structure
	subDir1 := filepath.Join(baseDir, "repo1")
	subDir2 := filepath.Join(baseDir, "repo2")
	nonRepoDir := filepath.Join(baseDir, "notrepo")

	os.Mkdir(subDir1, 0755)
	os.Mkdir(subDir2, 0755)
	os.Mkdir(nonRepoDir, 0755)
	os.Mkdir(filepath.Join(subDir1, ".git"), 0755)
	os.Mkdir(filepath.Join(subDir2, ".git"), 0755)

	httpSync := &repository.HTTPSync{}
	fileScanner := &repository.FileScanner{}
	storageManager := &repository.StorageManager{}
	scheduler := &service.Scheduler{}
	h := NewGRPCHandler(httpSync, fileScanner, storageManager, scheduler, mockLogger)

	// Test finding codebase paths
	configs, err := h.findCodebasePaths(baseDir, "test-name")
	assert.NoError(t, err)
	assert.Len(t, configs, 2) // Should find two git repositories

	// Verify returned configurations
	for _, config := range configs {
		switch config.CodebaseName {
		case "repo1":
			assert.Equal(t, subDir1, config.CodebasePath)
		case "repo2":
			assert.Equal(t, subDir2, config.CodebasePath)
		}
	}

	// Test invalid path
	_, err = h.findCodebasePaths(filepath.Join(baseDir, "nonexistent"), "test-name")
	assert.Error(t, err)
}
