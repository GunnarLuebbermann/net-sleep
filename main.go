package main

import (
	"fmt"
	"net-sleep/application"
	"net-sleep/structs"
)

func main() {
	cfg := structs.DefaultConfig()

	fmt.Println("🌙 NetSleep gestartet...")
	fmt.Printf("Überwache Netzwerktraffic alle %.0f Sekunden...\n", cfg.CheckInterval.Seconds())

	err := application.StartAutoShutdown(cfg)
	if err != nil {
		fmt.Println("Fehler:", err)
	}
}
