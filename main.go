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

	// ログをファイルに出力（デスクトップ or ホームディレクトリ）
	homeDir, _ := os.UserHomeDir()
	logPath := filepath.Join(homeDir, "Desktop", "razer_debug.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		// デスクトップがなければホームに書く
		logPath = filepath.Join(homeDir, "razer_debug.log")
		logFile, _ = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	}
	if logFile != nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}
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
