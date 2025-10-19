package structs

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test CheckInterval
	expectedCheckInterval := 10 * time.Second
	if config.CheckInterval != expectedCheckInterval {
		t.Errorf("Expected CheckInterval to be %v, got %v", expectedCheckInterval, config.CheckInterval)
	}

	// Test IdleThresholdBytes
	expectedIdleThresholdBytes := uint64(100 * 1024)
	if config.IdleThresholdBytes != expectedIdleThresholdBytes {
		t.Errorf("Expected IdleThresholdBytes to be %v, got %v", expectedIdleThresholdBytes, config.IdleThresholdBytes)
	}

	// Test IdleTimeBeforeAction
	expectedIdleTimeBeforeAction := 5 * time.Minute
	if config.IdleTimeBeforeAction != expectedIdleTimeBeforeAction {
		t.Errorf("Expected IdleTimeBeforeAction to be %v, got %v", expectedIdleTimeBeforeAction, config.IdleTimeBeforeAction)
	}

	// Test ShutdownCommand
	expectedShutdownCommand := "shutdown /s /t 0"
	if config.ShutdownCommand != expectedShutdownCommand {
		t.Errorf("Expected ShutdownCommand to be %q, got %q", expectedShutdownCommand, config.ShutdownCommand)
	}
}

func TestLoadConfig(t *testing.T) {
	// Test loading non-existent file (should return defaults)
	t.Run("NonExistentFile", func(t *testing.T) {
		config, err := LoadConfig("non_existent_file.json")
		if err != nil {
			t.Errorf("Expected no error for non-existent file, got %v", err)
		}

		defaultConfig := DefaultConfig()
		if config.CheckInterval != defaultConfig.CheckInterval {
			t.Errorf("Expected default CheckInterval %v, got %v", defaultConfig.CheckInterval, config.CheckInterval)
		}
	})

	// Test loading valid JSON file
	t.Run("ValidJSONFile", func(t *testing.T) {
		// Create temporary file
		content := `{
			"CheckIntervalSeconds": 30,
			"IdleThresholdBytes": 204800,
			"IdleTimeBeforeActionSeconds": 600,
			"ShutdownCommand": "poweroff"
		}`

		tmpFile := t.TempDir() + "/config.json"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		config, err := LoadConfig(tmpFile)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if config.CheckInterval != 30*time.Second {
			t.Errorf("Expected CheckInterval 30s, got %v", config.CheckInterval)
		}
		if config.IdleThresholdBytes != 204800 {
			t.Errorf("Expected IdleThresholdBytes 204800, got %v", config.IdleThresholdBytes)
		}
		if config.IdleTimeBeforeAction != 600*time.Second {
			t.Errorf("Expected IdleTimeBeforeAction 600s, got %v", config.IdleTimeBeforeAction)
		}
		if config.ShutdownCommand != "poweroff" {
			t.Errorf("Expected ShutdownCommand 'poweroff', got %q", config.ShutdownCommand)
		}
	})

	// Test partial configuration (should use defaults for missing fields)
	t.Run("PartialConfig", func(t *testing.T) {
		content := `{
			"CheckIntervalSeconds": 15,
			"ShutdownCommand": "halt"
		}`

		tmpFile := t.TempDir() + "/partial_config.json"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		config, err := LoadConfig(tmpFile)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		defaultConfig := DefaultConfig()
		if config.CheckInterval != 15*time.Second {
			t.Errorf("Expected CheckInterval 15s, got %v", config.CheckInterval)
		}
		if config.IdleThresholdBytes != defaultConfig.IdleThresholdBytes {
			t.Errorf("Expected default IdleThresholdBytes, got %v", config.IdleThresholdBytes)
		}
		if config.ShutdownCommand != "halt" {
			t.Errorf("Expected ShutdownCommand 'halt', got %q", config.ShutdownCommand)
		}
	})

	// Test invalid JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		content := `{
			"CheckIntervalSeconds": 30,
			"invalid": json
		}`

		tmpFile := t.TempDir() + "/invalid_config.json"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		_, err := LoadConfig(tmpFile)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	// Test zero values (should use defaults)
	t.Run("ZeroValues", func(t *testing.T) {
		content := `{
			"CheckIntervalSeconds": 0,
			"IdleThresholdBytes": 0,
			"IdleTimeBeforeActionSeconds": 0,
			"ShutdownCommand": ""
		}`

		tmpFile := t.TempDir() + "/zero_config.json"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		config, err := LoadConfig(tmpFile)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		defaultConfig := DefaultConfig()
		if config.CheckInterval != defaultConfig.CheckInterval {
			t.Errorf("Expected default CheckInterval, got %v", config.CheckInterval)
		}
		if config.IdleThresholdBytes != defaultConfig.IdleThresholdBytes {
			t.Errorf("Expected default IdleThresholdBytes, got %v", config.IdleThresholdBytes)
		}
		if config.ShutdownCommand != defaultConfig.ShutdownCommand {
			t.Errorf("Expected default ShutdownCommand, got %q", config.ShutdownCommand)
		}
	})
}
