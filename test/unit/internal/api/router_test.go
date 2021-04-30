package api_test

import (
	"testing"
	"uacl/internal/api"
)

func TestRouterInitializes(t *testing.T) {
	if err := api.CreateRouter(); err == nil {
		t.Errorf("Expected router to not be nil")
	}
}
