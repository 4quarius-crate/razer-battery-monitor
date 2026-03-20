package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-toast/toast"
)

// Tray はシステムトレイアイコンと通知を管理する
type Tray struct {
	cfg        Config
	monitor    *Monitor
	lastNotify time.Time
}

func NewTray(cfg Config, monitor *Monitor) *Tray {
	return &Tray{cfg: cfg, monitor: monitor}
}

// Run はトレイアイコンを起動する（メインスレッドをブロック）
func (t *Tray) Run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *Tray) onReady() {
	systray.SetIcon(makeBatteryIcon(100, false))
	systray.SetTooltip("Razer Battery Monitor — 接続待ち")

	mRefresh := systray.AddMenuItem("今すぐ更新", "バッテリー残量を今すぐ確認")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("終了", "Razer Battery Monitor を終了")

	go func() {
		for {
			select {
			case <-mRefresh.ClickedCh:
				t.monitor.Refresh()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.monitor.Stop()
	log.Println("終了しました")
}

// Update はバッテリー状態を受け取りアイコンとツールチップを更新する
func (t *Tray) Update(status BatteryStatus) {
	systray.SetIcon(makeBatteryIcon(status.Level, status.IsCharging))

	chargeLabel := ""
	if status.IsCharging {
		chargeLabel = " (充電中)"
	}
	systray.SetTooltip(fmt.Sprintf("Razer マウス: %d%%%s", status.Level, chargeLabel))

	// 低バッテリー通知
	cooldown := time.Duration(t.cfg.Alert.NotifyCooldown) * time.Second
	if !status.IsCharging &&
		status.Level <= t.cfg.Alert.LowBatteryThreshold &&
		time.Since(t.lastNotify) > cooldown {
		t.showToast(status.Level)
		t.lastNotify = time.Now()
	}
}

func (t *Tray) showToast(level int) {
	notification := toast.Notification{
		AppID:   "Razer Battery Monitor",
		Title:   "バッテリー残量が低下しています",
		Message: fmt.Sprintf("Razer マウスのバッテリーが %d%% です", level),
	}
	if err := notification.Push(); err != nil {
		log.Printf("通知の表示に失敗: %v", err)
	}
}
