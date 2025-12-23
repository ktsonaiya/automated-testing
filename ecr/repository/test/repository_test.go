package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Repository represents an ECR repository
type Repository struct {
	Name                      string
	ImageScanningConfiguration bool
	ImageTagMutability        string
	EncryptionType            string
}

// NewRepository creates a new ECR Repository instance
func NewRepository(name string) *Repository {
	return &Repository{
		Name:                      name,
		ImageScanningConfiguration: false,
		ImageTagMutability:        "MUTABLE",
		EncryptionType:            "AES256",
	}
}

// EnableImageScanning enables image scanning for the repository
func (r *Repository) EnableImageScanning() {
	r.ImageScanningConfiguration = true
}

// SetImageTagMutability sets the image tag mutability
func (r *Repository) SetImageTagMutability(mutability string) bool {
	if mutability != "MUTABLE" && mutability != "IMMUTABLE" {
		return false
	}
	r.ImageTagMutability = mutability
	return true
}

// IsValidRepository validates the repository configuration
func (r *Repository) IsValidRepository() bool {
	validMutability := r.ImageTagMutability == "MUTABLE" || r.ImageTagMutability == "IMMUTABLE"
	validEncryption := r.EncryptionType == "AES256" || r.EncryptionType == "KMS"
	return len(r.Name) > 0 && validMutability && validEncryption
}

// TestCreateRepository tests the creation of an ECR repository
func TestCreateRepository(t *testing.T) {
	repo := NewRepository("test-repo")

	assert.Equal(t, "test-repo", repo.Name)
	assert.False(t, repo.ImageScanningConfiguration)
	assert.Equal(t, "MUTABLE", repo.ImageTagMutability)
	assert.Equal(t, "AES256", repo.EncryptionType)
}

// TestEnableImageScanning tests enabling image scanning
func TestEnableImageScanning(t *testing.T) {
	repo := NewRepository("test-repo")
	assert.False(t, repo.ImageScanningConfiguration)

	repo.EnableImageScanning()
	assert.True(t, repo.ImageScanningConfiguration)
}

// TestSetImageTagMutability tests setting image tag mutability
func TestSetImageTagMutability(t *testing.T) {
	repo := NewRepository("test-repo")

	tests := []struct {
		name        string
		mutability  string
		shouldSucceed bool
	}{
		{"valid - MUTABLE", "MUTABLE", true},
		{"valid - IMMUTABLE", "IMMUTABLE", true},
		{"invalid - INVALID", "INVALID", false},
		{"invalid - empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := repo.SetImageTagMutability(tt.mutability)
			assert.Equal(t, tt.shouldSucceed, success)
			if tt.shouldSucceed {
				assert.Equal(t, tt.mutability, repo.ImageTagMutability)
			}
		})
	}
}

// TestRepositoryValidation tests repository validation
func TestRepositoryValidation(t *testing.T) {
	tests := []struct {
		name      string
		repo      *Repository
		isValid   bool
	}{
		{
			name: "valid repository",
			repo: &Repository{
				Name:               "valid-repo",
				ImageTagMutability: "MUTABLE",
				EncryptionType:     "AES256",
			},
			isValid: true,
		},
		{
			name: "valid repository with KMS encryption",
			repo: &Repository{
				Name:               "kms-repo",
				ImageTagMutability: "IMMUTABLE",
				EncryptionType:     "KMS",
			},
			isValid: true,
		},
		{
			name: "invalid - empty name",
			repo: &Repository{
				Name:               "",
				ImageTagMutability: "MUTABLE",
				EncryptionType:     "AES256",
			},
			isValid: false,
		},
		{
			name: "invalid - invalid encryption type",
			repo: &Repository{
				Name:               "invalid-repo",
				ImageTagMutability: "MUTABLE",
				EncryptionType:     "INVALID",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.repo.IsValidRepository())
		})
	}
}
