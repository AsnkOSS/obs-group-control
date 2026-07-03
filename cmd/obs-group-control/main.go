// Command obs-group-control starts/stops recording on a group of OBS
// instances listed in devices.ini.
package main

import (
	"flag"
	"log"

	"obs-group-control/internal/config"
	"obs-group-control/internal/obs"
	"obs-group-control/internal/ui"
)

func main() {
	configPath := flag.String("config", "devices.ini", "path to device list")
	flag.Parse()

	devices, err := config.LoadDevices(*configPath)
	if err != nil {
		log.Fatalf("load devices: %v", err)
	}

	controller := obs.NewController(devices)
	ui.New(controller).Run()
}
