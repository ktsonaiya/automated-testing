package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AutoScalingGroup represents a simplified EC2 Auto Scaling Group
type AutoScalingGroup struct {
	Name                   string
	MinSize                int
	MaxSize                int
	DesiredCapacity        int
	HealthCheckType        string
	HealthCheckGracePeriod int
}

// NewAutoScalingGroup creates a new Auto Scaling Group instance
func NewAutoScalingGroup(name string, minSize, maxSize, desiredCapacity int) *AutoScalingGroup {
	return &AutoScalingGroup{
		Name:                   name,
		MinSize:                minSize,
		MaxSize:                maxSize,
		DesiredCapacity:        desiredCapacity,
		HealthCheckType:        "ELB",
		HealthCheckGracePeriod: 300,
	}
}

// IsValidConfiguration validates the Auto Scaling Group configuration
func (asg *AutoScalingGroup) IsValidConfiguration() bool {
	return asg.MinSize >= 0 &&
		asg.MaxSize > asg.MinSize &&
		asg.DesiredCapacity >= asg.MinSize &&
		asg.DesiredCapacity <= asg.MaxSize &&
		asg.HealthCheckGracePeriod > 0
}

// TestCreateAutoScalingGroup tests the creation of an Auto Scaling Group
func TestCreateAutoScalingGroup(t *testing.T) {
	asg := NewAutoScalingGroup("test-asg", 2, 10, 5)

	assert.Equal(t, "test-asg", asg.Name)
	assert.Equal(t, 2, asg.MinSize)
	assert.Equal(t, 10, asg.MaxSize)
	assert.Equal(t, 5, asg.DesiredCapacity)
	assert.Equal(t, "ELB", asg.HealthCheckType)
	assert.Equal(t, 300, asg.HealthCheckGracePeriod)
}

// TestAutoScalingGroupValidation tests the validation of Auto Scaling Group configuration
func TestAutoScalingGroupValidation(t *testing.T) {
	tests := []struct {
		name    string
		asg     *AutoScalingGroup
		isValid bool
	}{
		{
			name:    "valid configuration",
			asg:     NewAutoScalingGroup("valid-asg", 1, 10, 5),
			isValid: true,
		},
		{
			name: "invalid - desired capacity below minimum",
			asg: &AutoScalingGroup{
				Name:                   "invalid-asg",
				MinSize:                5,
				MaxSize:                10,
				DesiredCapacity:        2,
				HealthCheckGracePeriod: 300,
			},
			isValid: false,
		},
		{
			name: "invalid - desired capacity above maximum",
			asg: &AutoScalingGroup{
				Name:                   "invalid-asg",
				MinSize:                1,
				MaxSize:                5,
				DesiredCapacity:        10,
				HealthCheckGracePeriod: 300,
			},
			isValid: false,
		},
		{
			name: "invalid - max size not greater than min size",
			asg: &AutoScalingGroup{
				Name:                   "invalid-asg",
				MinSize:                10,
				MaxSize:                5,
				DesiredCapacity:        7,
				HealthCheckGracePeriod: 300,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.asg.IsValidConfiguration())
		})
	}
}

// TestAutoScalingGroupDefaults tests default values
func TestAutoScalingGroupDefaults(t *testing.T) {
	asg := NewAutoScalingGroup("default-test", 1, 5, 3)

	assert.Equal(t, "ELB", asg.HealthCheckType)
	assert.Equal(t, 300, asg.HealthCheckGracePeriod)
}
