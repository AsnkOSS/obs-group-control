// Package config loads the device list from an INI file.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Device is one OBS instance endpoint.
type Device struct {
	Name     string // section name from the INI file
	Addr     string // host:port
	Password string // optional websocket password
}

// LoadDevices parses an INI file where each section is one device:
//
//	[cam1]
//	address = 192.168.1.10:4455
//	password = secret   ; optional
//
// Bare "host:port" lines outside any section are also accepted,
// so a minimal file can be a single address line without a password.
func LoadDevices(path string) ([]Device, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		devices []Device
		current *Device
	)
	flush := func() error {
		if current == nil {
			return nil
		}
		if current.Addr == "" {
			return fmt.Errorf("device %q: missing address", current.Name)
		}
		devices = append(devices, *current)
		current = nil
		return nil
	}

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := stripComment(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]"):
			if err := flush(); err != nil {
				return nil, err
			}
			current = &Device{Name: strings.TrimSpace(line[1 : len(line)-1])}

		case strings.Contains(line, "="):
			key, value, _ := strings.Cut(line, "=")
			key = strings.ToLower(strings.TrimSpace(key))
			value = strings.TrimSpace(value)
			if current == nil {
				return nil, fmt.Errorf("line %d: key %q outside a [section]", lineNo, key)
			}
			switch key {
			case "address", "addr", "host":
				current.Addr = value
			case "password":
				current.Password = value
			default:
				return nil, fmt.Errorf("line %d: unknown key %q", lineNo, key)
			}

		default:
			// Bare address line, e.g. "127.0.0.1:4455".
			if err := flush(); err != nil {
				return nil, err
			}
			devices = append(devices, Device{Name: line, Addr: line})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := flush(); err != nil {
		return nil, err
	}

	for _, d := range devices {
		if !strings.Contains(d.Addr, ":") {
			return nil, fmt.Errorf("device %q: address %q missing port", d.Name, d.Addr)
		}
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found in %s", path)
	}
	return devices, nil
}

// stripComment trims whitespace and removes ';' or '#' comments.
func stripComment(line string) string {
	if i := strings.IndexAny(line, ";#"); i >= 0 {
		line = line[:i]
	}
	return strings.TrimSpace(line)
}
