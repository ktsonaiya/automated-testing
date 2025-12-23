package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Service represents an ECS service
type Service struct {
	Name            string
	Cluster         string
	DesiredCount    int
	LaunchType      string
	NetworkMode     string
	HealthCheckPath string
	Enabled         bool
}

// NewService creates a new ECS Service instance
func NewService(name, cluster string, desiredCount int) *Service {
	return &Service{
		Name:            name,
		Cluster:         cluster,
		DesiredCount:    desiredCount,
		LaunchType:      "FARGATE",
		NetworkMode:     "awsvpc",
		HealthCheckPath: "/health",
		Enabled:         true,
	}
}

// Scale updates the desired count of the service
func (s *Service) Scale(count int) bool {
	if count < 0 {
		return false
	}
	s.DesiredCount = count
	return true
}

// SetLaunchType sets the launch type for the service
func (s *Service) SetLaunchType(launchType string) bool {
	if launchType != "FARGATE" && launchType != "EC2" && launchType != "EXTERNAL" {
		return false
	}
	s.LaunchType = launchType
	return true
}

// IsValidService validates the service configuration
func (s *Service) IsValidService() bool {
	validLaunchType := s.LaunchType == "FARGATE" || s.LaunchType == "EC2" || s.LaunchType == "EXTERNAL"
	validNetworkMode := s.NetworkMode == "awsvpc" || s.NetworkMode == "bridge" || s.NetworkMode == "host" || s.NetworkMode == "none"
	return len(s.Name) > 0 &&
		len(s.Cluster) > 0 &&
		s.DesiredCount >= 0 &&
		validLaunchType &&
		validNetworkMode &&
		s.Enabled
}

// TestCreateService tests the creation of an ECS service
func TestCreateService(t *testing.T) {
	service := NewService("test-service", "test-cluster", 3)

	assert.Equal(t, "test-service", service.Name)
	assert.Equal(t, "test-cluster", service.Cluster)
	assert.Equal(t, 3, service.DesiredCount)
	assert.Equal(t, "FARGATE", service.LaunchType)
	assert.Equal(t, "awsvpc", service.NetworkMode)
	assert.Equal(t, "/health", service.HealthCheckPath)
	assert.True(t, service.Enabled)
}

// TestScaleService tests scaling a service
func TestScaleService(t *testing.T) {
	service := NewService("test-service", "test-cluster", 3)

	tests := []struct {
		name         string
		count        int
		shouldSucceed bool
		expectedCount int
	}{
		{"scale up", 5, true, 5},
		{"scale down", 1, true, 1},
		{"scale to zero", 0, true, 0},
		{"invalid - negative count", -1, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCount := service.DesiredCount
			success := service.Scale(tt.count)
			assert.Equal(t, tt.shouldSucceed, success)
			if success {
				assert.Equal(t, tt.expectedCount, service.DesiredCount)
			} else {
				assert.Equal(t, originalCount, service.DesiredCount)
			}
		})
	}
}

// TestSetLaunchType tests setting the launch type
func TestSetLaunchType(t *testing.T) {
	service := NewService("test-service", "test-cluster", 3)

	tests := []struct {
		name         string
		launchType   string
		shouldSucceed bool
	}{
		{"valid - FARGATE", "FARGATE", true},
		{"valid - EC2", "EC2", true},
		{"valid - EXTERNAL", "EXTERNAL", true},
		{"invalid - INVALID", "INVALID", false},
		{"invalid - empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := service.SetLaunchType(tt.launchType)
			assert.Equal(t, tt.shouldSucceed, success)
			if success {
				assert.Equal(t, tt.launchType, service.LaunchType)
			}
		})
	}
}

// TestServiceValidation tests service validation
func TestServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service *Service
		isValid bool
	}{
		{
			name: "valid FARGATE service",
			service: &Service{
				Name:        "valid-service",
				Cluster:     "valid-cluster",
				DesiredCount: 2,
				LaunchType:  "FARGATE",
				NetworkMode: "awsvpc",
				Enabled:     true,
			},
			isValid: true,
		},
		{
			name: "valid EC2 service",
			service: &Service{
				Name:        "valid-service",
				Cluster:     "valid-cluster",
				DesiredCount: 1,
				LaunchType:  "EC2",
				NetworkMode: "bridge",
				Enabled:     true,
			},
			isValid: true,
		},
		{
			name: "invalid - empty name",
			service: &Service{
				Name:        "",
				Cluster:     "valid-cluster",
				DesiredCount: 2,
				LaunchType:  "FARGATE",
				NetworkMode: "awsvpc",
				Enabled:     true,
			},
			isValid: false,
		},
		{
			name: "invalid - empty cluster",
			service: &Service{
				Name:        "valid-service",
				Cluster:     "",
				DesiredCount: 2,
				LaunchType:  "FARGATE",
				NetworkMode: "awsvpc",
				Enabled:     true,
			},
			isValid: false,
		},
		{
			name: "invalid - disabled service",
			service: &Service{
				Name:        "valid-service",
				Cluster:     "valid-cluster",
				DesiredCount: 2,
				LaunchType:  "FARGATE",
				NetworkMode: "awsvpc",
				Enabled:     false,
			},
			isValid: false,
		},
		{
			name: "invalid - negative desired count",
			service: &Service{
				Name:        "valid-service",
				Cluster:     "valid-cluster",
				DesiredCount: -1,
				LaunchType:  "FARGATE",
				NetworkMode: "awsvpc",
				Enabled:     true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.service.IsValidService())
		})
	}
}
