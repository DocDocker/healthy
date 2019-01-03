package checker

import (
	"github.com/localghost/healthy/utils"
	"testing"
)

func TestFailureOnParsing(t *testing.T) {
	checker, err := NewChecker(42)
	if err == nil {
		t.Fatal("Expected creating checker to fail with error")
	}
	if checker != nil {
		t.Fatal("Expected creating checker to return nil")
	}
}

func TestCreatingEmptyChecker(t *testing.T) {
	checker, err := NewChecker(map[string]interface{}{})
	if err != nil {
		t.Fatal("Expected creating empty checker to succeed")
	}
	if checker == nil {
		t.Fatal("Expected creating empty checker to return checker object")
	}
}

func TestGetMissingCheck(t *testing.T) {
	checker, _ := NewChecker(map[string]interface{}{})
	if err, ok := checker.Get("foo").(utils.NoSuchCheckError); !ok {
		t.Fatalf("Expected NoSuchCheckError but got %#v", err)
	}
}

func TestGet(t *testing.T) {
	checker, _ := NewChecker(map[string]interface{}{
		"ls": map[string]interface{} {
			"type": "command",
			"command": "ls",
		},
	})
	if err := checker.Get("ls"); err != nil {
		t.Fatalf("Expected check to succeed but it failed with %#v", err)
	}
}

func TestGetAll(t *testing.T) {
	checker, _ := NewChecker(map[string]interface{}{
		"ls": map[string]interface{} {
			"type": "command",
			"command": "ls",
		},
		"echo": map[string]interface{} {
			"type": "command",
			"command": "echo",
		},
	})
	if err := checker.GetAll(); err != nil {
		t.Fatalf("Expected check to succeed but it failed with %#v", err)
	}
}
