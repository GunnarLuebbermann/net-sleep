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
//
// Fields:
//   - CheckIntervalSeconds: How often (in seconds) to check network activity
//   - IdleThresholdBytes: Maximum bytes per second to consider the network idle
//   - IdleTimeBeforeActionSeconds: How long (in seconds) the network must be idle before triggering an action
//   - ShutdownCommand: The system command to execute when idle threshold is met
type ConfigJSON struct {
	CheckIntervalSeconds        int    `json:"CheckIntervalSeconds"`
	IdleThresholdBytes          uint64 `json:"IdleThresholdBytes"`
	IdleTimeBeforeActionSeconds int    `json:"IdleTimeBeforeActionSeconds"`
	ShutdownCommand             string `json:"ShutdownCommand"`
}

// Config holds the configuration parameters for the network activity monitoring system.
// It defines thresholds, intervals, and actions for monitoring network traffic
// and triggering system shutdown when the network is idle for a specified duration.
type Config struct {
	CheckInterval        time.Duration // Zeit zwischen Messungen
	IdleThresholdBytes   uint64        // Grenzwert in Bytes/s (unterhalb = inaktiv)
	IdleTimeBeforeAction time.Duration // Wie lange inaktiv, bevor Aktion ausgelöst wird
	ShutdownCommand      string        // Systembefehl zum Herunterfahren
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
		IdleTimeBeforeAction: 1 * time.Minute,
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
	var cfg Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, fmt.Errorf("fehler beim Lesen der Datei: %v", err)
	}

	var cfgJSON ConfigJSON
	if err := json.Unmarshal(data, &cfgJSON); err != nil {
		return cfg, fmt.Errorf("fehler beim Parsen der JSON: %v", err)
	}

	// Validiere und konvertiere Konfigurationswerte
	if cfgJSON.CheckIntervalSeconds <= 0 {
		return cfg, fmt.Errorf("CheckIntervalSeconds muss größer als 0 sein")
	}
	if cfgJSON.IdleTimeBeforeActionSeconds <= 0 {
		return cfg, fmt.Errorf("IdleTimeBeforeActionSeconds muss größer als 0 sein")
	}
	if cfgJSON.ShutdownCommand == "" {
		return cfg, fmt.Errorf("ShutdownCommand darf nicht leer sein")
	}

	// Konvertiere Sekunden in time.Duration
	cfg.CheckInterval = time.Duration(cfgJSON.CheckIntervalSeconds) * time.Second
	cfg.IdleThresholdBytes = cfgJSON.IdleThresholdBytes
	cfg.IdleTimeBeforeAction = time.Duration(cfgJSON.IdleTimeBeforeActionSeconds) * time.Second
	cfg.ShutdownCommand = cfgJSON.ShutdownCommand

	return cfg, nil
}
