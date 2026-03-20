package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Monitor  MonitorConfig  `toml:"monitor"`
	Alert    AlertConfig    `toml:"alert"`
	FPSGuard FPSGuardConfig `toml:"fps_guard"`
	Debug    DebugConfig    `toml:"debug"`
}

type MonitorConfig struct {
	PollIntervalNormal   int `toml:"poll_interval_normal"`
	PollIntervalCharging int `toml:"poll_interval_charging"`
	PollIntervalLow      int `toml:"poll_interval_low"`
}

type AlertConfig struct {
	LowBatteryThreshold int `toml:"low_battery_threshold"`
	NotifyCooldown      int `toml:"notify_cooldown"`
}

type FPSGuardConfig struct {
	Enabled       bool     `toml:"enabled"`
	GameProcesses []string `toml:"game_processes"`
}

type DebugConfig struct {
	DumpRawBytes bool `toml:"dump_raw_bytes"`
}

func defaultConfig() Config {
	return Config{
		Monitor: MonitorConfig{
			PollIntervalNormal:   300,
			PollIntervalCharging: 600,
			PollIntervalLow:      120,
		},
		Alert: AlertConfig{
			LowBatteryThreshold: 20,
			NotifyCooldown:      1800,
		},
		FPSGuard: FPSGuardConfig{
			Enabled: true,
		},
	}
}

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil // ファイルがなければデフォルト値を使用
	}
	_, err = toml.Decode(string(data), &cfg)
	return cfg, err
}
