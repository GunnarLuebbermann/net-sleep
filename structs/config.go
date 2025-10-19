package structs

import "time"

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
// - IdleThresholdBytes to 50KB (minimum network activity to consider as "active")
// - IdleTimeBeforeAction to 5 minutes (duration of inactivity before triggering action)
// - ShutdownCommand to "shutdown /s /t 0" (Windows shutdown command)
func DefaultConfig() Config {
	return Config{
		CheckInterval:        10 * time.Second,
		IdleThresholdBytes:   50 * 1024,
		IdleTimeBeforeAction: 5 * time.Minute,
		ShutdownCommand:      "shutdown /s /t 0",
	}
}
