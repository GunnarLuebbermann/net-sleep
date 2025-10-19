package main

import (
	"fmt"
	"net-sleep/application"
	"net-sleep/structs"
)

func main() {
	cfg, err := structs.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Keine gültige Konfiguration gefunden. Benutze Default Konfiguration.")
		cfg = structs.DefaultConfig()
	}

	fmt.Println("🌙 NetSleep gestartet...")
	fmt.Println("📊 Überwache Netzwerktraffic...")

	fmt.Printf("Verwendete Konfiguration:\n")
	fmt.Printf("- Check Interval: %.0f Sekunden\n", cfg.CheckInterval.Seconds())
	fmt.Printf("- Traffic Threshold: %s\n", application.FormatBytes(cfg.IdleThresholdBytes))
	fmt.Printf("- Shutdown Delay: %.0f Sekunden\n", cfg.IdleTimeBeforeAction.Seconds())
	fmt.Printf("- Shutdown Command: %s\n", cfg.ShutdownCommand)
	fmt.Println()

	err = application.StartAutoShutdown(cfg)
	if err != nil {
		fmt.Println("Fehler:", err)
	}
}
