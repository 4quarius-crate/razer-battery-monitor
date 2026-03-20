package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	// FPS干渉防止: プロセス優先度を下げる
	setLowPriority()

	// 設定読み込み（exe と同じフォルダの data/config.toml）
	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configPath := filepath.Join(exeDir, "data", "config.toml")
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Printf("設定ファイルの読み込み失敗: %v (デフォルト値を使用)", err)
	}

	// 監視とトレイを初期化
	var tray *Tray
	monitor := NewMonitor(cfg, func(status BatteryStatus) {
		if tray != nil {
			tray.Update(status)
		}
	})
	tray = NewTray(cfg, monitor)

	monitor.Start()

	// トレイアイコン起動（終了まで待機）
	tray.Run()
}
