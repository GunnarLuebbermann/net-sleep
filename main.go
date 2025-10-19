package main

import (
	"fmt"
	"net-sleep/application"
	"net-sleep/structs"
)

func main() {
	cfg, err := structs.LoadConfig("config.json")
	if err != nil {
		fmt.Println("No valid configuration found. Using default configuration.")
		cfg = structs.DefaultConfig()
	}

	fmt.Println("🌙 NetSleep started...")
	fmt.Println("📊 Monitoring network traffic...")

	fmt.Printf("Configuration used:\n")
	fmt.Printf("- Check Interval: %.0f seconds\n", cfg.CheckInterval.Seconds())
	fmt.Printf("- Traffic Threshold: %s\n", application.FormatBytes(cfg.IdleThresholdBytes))
	fmt.Printf("- Shutdown Delay: %.0f seconds\n", cfg.IdleTimeBeforeAction.Seconds())
	fmt.Printf("- Shutdown Command: %s\n", cfg.ShutdownCommand)
	fmt.Println()

	err = application.StartAutoShutdown(cfg)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
