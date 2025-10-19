package application

import (
	"runtime"
	"strings"
	"testing"
)

func TestShutdownPC(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "valid echo command",
			command:     "echo test",
			expectError: false,
		},
		{
			name:        "invalid command",
			command:     "nonexistentcommand12345",
			expectError: true,
		},
		{
			name:        "empty command",
			command:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests on non-Windows systems since the function is Windows-specific
			if runtime.GOOS != "windows" {
				t.Skip("shutdownPC is Windows-specific")
			}

			err := shutdownPC(tt.command)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestShutdownPCCommand(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("shutdownPC is Windows-specific")
	}

	// Test that the function accepts typical shutdown commands
	commands := []string{
		"shutdown /s /t 0",
		"shutdown /r /t 0",
		"shutdown /l",
	}

	for _, cmd := range commands {
		t.Run("command_"+strings.ReplaceAll(cmd, " ", "_"), func(t *testing.T) {
			// We don't actually execute shutdown commands in tests
			// Instead, we test with harmless commands that use similar syntax
			testCmd := "echo " + cmd
			err := shutdownPC(testCmd)
			if err != nil {
				t.Errorf("failed to execute test command: %v", err)
			}
		})
	}
}

func TestGetNetworkBytes(t *testing.T) {
	t.Run("successful network bytes retrieval", func(t *testing.T) {
		sent, recv, err := getNetworkBytes()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Network bytes should be non-negative
		if sent < 0 {
			t.Errorf("bytes sent should be non-negative, got %d", sent)
		}

		if recv < 0 {
			t.Errorf("bytes received should be non-negative, got %d", recv)
		}

		// In a typical system, there should be some network activity
		// but we can't guarantee specific values, so we just check they're valid uint64
		t.Logf("Bytes sent: %d, Bytes received: %d", sent, recv)
	})

	t.Run("return values are consistent", func(t *testing.T) {
		// Call the function twice and ensure it returns valid data both times
		sent1, recv1, err1 := getNetworkBytes()
		if err1 != nil {
			t.Fatalf("first call failed: %v", err1)
		}

		sent2, recv2, err2 := getNetworkBytes()
		if err2 != nil {
			t.Fatalf("second call failed: %v", err2)
		}

		// Second call should have equal or greater values (monotonic increase)
		if sent2 < sent1 {
			t.Errorf("bytes sent decreased between calls: %d -> %d", sent1, sent2)
		}

		if recv2 < recv1 {
			t.Errorf("bytes received decreased between calls: %d -> %d", recv1, recv2)
		}
	})
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "single byte",
			bytes:    1,
			expected: "1 B",
		},
		{
			name:     "bytes under 1KB",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "exactly 1KB",
			bytes:    1024,
			expected: "1.00 KB",
		},
		{
			name:     "1.5KB",
			bytes:    1536,
			expected: "1.50 KB",
		},
		{
			name:     "bytes under 1MB",
			bytes:    512 * 1024,
			expected: "512.00 KB",
		},
		{
			name:     "exactly 1MB",
			bytes:    1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "2MB",
			bytes:    2 * 1024 * 1024,
			expected: "2.00 MB",
		},
		{
			name:     "1.5MB",
			bytes:    1536 * 1024,
			expected: "1.50 MB",
		},
		{
			name:     "bytes under 1GB",
			bytes:    512 * 1024 * 1024,
			expected: "512.00 MB",
		},
		{
			name:     "exactly 1GB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "2.5GB",
			bytes:    2560 * 1024 * 1024,
			expected: "2.50 GB",
		},
		{
			name:     "large value",
			bytes:    5 * 1024 * 1024 * 1024,
			expected: "5.00 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		})
	}
}
