package application

import (
	"fmt"
	"net-sleep/structs"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// StartAutoShutdown monitors network traffic and initiates system shutdown when network activity
// falls below the configured threshold for a specified duration.
//
// The function continuously polls network statistics at intervals defined by cfg.CheckInterval
// and calculates the network traffic speed in bytes per second. If the traffic remains below
// cfg.IdleThresholdBytes for a period equal to or greater than cfg.IdleTimeBeforeAction,
// the system shutdown process is triggered.
//
// Parameters:
//   - cfg: Configuration struct containing monitoring parameters including:
//   - CheckInterval: Duration between network traffic checks
//   - IdleThresholdBytes: Minimum bytes/second threshold to consider as active
//   - IdleTimeBeforeAction: Duration of low activity before triggering shutdown
//   - ShutdownCommand: Command to execute for system shutdown (currently commented out)
//
// Returns:
//   - error: Returns an error if network data retrieval fails or shutdown command fails
//
// The function runs indefinitely until an error occurs or shutdown is initiated.
// Network statistics and idle duration are logged to stdout for monitoring purposes.
func StartAutoShutdown(cfg structs.Config) error {
	var lastSent, lastRecv uint64
	var idleDuration time.Duration

	fmt.Println("📡 Starting network monitoring...")

	for {
		sent, recv, err := getNetworkBytes()
		if err != nil {
			fmt.Println("❌ Error retrieving network data:", err)
			time.Sleep(cfg.CheckInterval)
			continue
		}

		if lastSent != 0 {
			sentDiff := sent - lastSent
			recvDiff := recv - lastRecv
			totalDiff := sentDiff + recvDiff
			speed := float64(totalDiff) / cfg.CheckInterval.Seconds()

			// Format speed with appropriate unit
			fmt.Printf("📊 Traffic: %s/s\n", FormatBytes(uint64(speed)))

			if speed < float64(cfg.IdleThresholdBytes) {
				idleDuration += cfg.CheckInterval
				fmt.Printf("⚠️  Low traffic for %v\n", idleDuration)
				if idleDuration >= cfg.IdleTimeBeforeAction {
					fmt.Println("🛑 No activity detected – shutting down PC...")
					return shutdownPC(cfg.ShutdownCommand)
				}
			} else {
				idleDuration = 0
			}
		}

		lastSent = sent
		lastRecv = recv
		time.Sleep(cfg.CheckInterval)
	}
}

// getNetworkBytes retrieves the total bytes sent and received across all network interfaces.
// It returns the bytes sent, bytes received, and any error encountered.
// If no network interfaces are found, it returns an error.
// The function uses the first available network interface from the system's IO counters.
func getNetworkBytes() (uint64, uint64, error) {
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}
	if len(ioCounters) == 0 {
		return 0, 0, fmt.Errorf("no network interfaces found")
	}
	return ioCounters[0].BytesSent, ioCounters[0].BytesRecv, nil
}

// FormatBytes converts a byte count into a human-readable string representation
// with appropriate units (B, KB, MB, GB). The function automatically selects
// the most appropriate unit based on the size of the input value and formats
// the result with two decimal places for units larger than bytes.
//
// Parameters:
//   - bytes: The number of bytes to format as a uint64 value
//
// Returns:
//   - A string representation of the byte count with appropriate unit suffix
//
// Examples:
//   - FormatBytes(512) returns "512 B"
//   - FormatBytes(1536) returns "1.50 KB"
//   - FormatBytes(2097152) returns "2.00 MB"
//   - FormatBytes(1073741824) returns "1.00 GB"
func FormatBytes(bytes uint64) string {
	if bytes >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
	} else if bytes >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
	} else if bytes >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%d B", bytes)
	}
}

// shutdownPC executes a system shutdown command on the current operating system.
// It takes a command string and executes it using the appropriate shell.
// On Windows, it uses cmd.exe with /C flag to execute the command.
// Returns an error if the command execution fails.
//
// Parameters:
//   - command: The shutdown command string to execute
//
// Returns:
//   - error: nil on successful execution, error otherwise
//
// Example:
//
//	err := shutdownPC("shutdown /s /t 0")
func shutdownPC(command string) error {
	fmt.Println("🚦 Executing shutdown command:", command)
	parts := []string{"cmd", "/C", command}
	// Linux: parts := []string{"bash", "-c", command}
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run()
}
