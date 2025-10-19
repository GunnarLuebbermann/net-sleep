package structs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ConfigJSON represents the configuration structure for network sleep monitoring.
// It defines the parameters used to monitor network activity and trigger actions
// when the system is idle for a specified duration.
type ConfigJSON struct {
	Interval        int    `json:"interval"`
	IdleThreshold   uint64 `json:"idle_threshold"`
	IdleTime        int    `json:"idle_time"`
	ShutdownCommand string `json:"shutdown_command"`
}

// Config holds configuration parameters for monitoring network activity and performing actions
// when the system is idle. It specifies the interval for checking activity, the threshold for
// considering the system idle based on network usage, the duration to wait before taking action,
// and the command to execute when the idle condition is met.
type Config struct {
	CheckInterval        time.Duration
	IdleThresholdBytes   uint64
	IdleTimeBeforeAction time.Duration
	ShutdownCommand      string
}

// DefaultConfig returns a Config struct with sensible default values for network monitoring.
// The default configuration sets:
// - CheckInterval to 10 seconds (frequency of network activity checks)
// - IdleThresholdBytes to 100KB (minimum network activity to consider as "active")
// - IdleTimeBeforeAction to 1 minute (duration of inactivity before triggering action)
// - ShutdownCommand to "shutdown /s /t 0" (Windows shutdown command)
func DefaultConfig() Config {
	return Config{
		CheckInterval:        10 * time.Second,
		IdleThresholdBytes:   100 * 1024,
		IdleTimeBeforeAction: 5 * time.Minute,
		ShutdownCommand:      "shutdown /s /t 0",
	}
}

// LoadConfig reads a configuration file from the specified filename and returns a Config struct.
// It parses the JSON data from the file and converts time values from seconds to time.Duration.
// Returns an error if the file cannot be read or if the JSON parsing fails.
//
// Parameters:
//   - filename: The path to the configuration file to load
//
// Returns:
//   - Config: The parsed configuration with converted time durations
//   - error: An error if file reading or JSON parsing fails
func LoadConfig(filename string) (Config, error) {
	cfg := DefaultConfig() // Start with defaults

	data, err := os.ReadFile(filename)
	if err != nil {
		// If file is missing, simply return defaults, no error
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("error reading file: %v", err)
	}

	var cfgJSON ConfigJSON
	if err := json.Unmarshal(data, &cfgJSON); err != nil {
		return cfg, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Only overwrite fields that are set
	if cfgJSON.Interval > 0 {
		cfg.CheckInterval = time.Duration(cfgJSON.Interval) * time.Second
	}
	if cfgJSON.IdleThreshold > 0 {
		cfg.IdleThresholdBytes = cfgJSON.IdleThreshold
	}
	if cfgJSON.IdleTime > 0 {
		cfg.IdleTimeBeforeAction = time.Duration(cfgJSON.IdleTime) * time.Second
	}
	if cfgJSON.ShutdownCommand != "" {
		cfg.ShutdownCommand = cfgJSON.ShutdownCommand
	}

	return cfg, nil
}
