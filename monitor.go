package main

import (
	"log"
	"time"
)

// Monitor はバッテリー残量をポーリングし、更新をコールバックに通知する
type Monitor struct {
	cfg        Config
	onUpdate   func(BatteryStatus)
	stop       chan struct{}
	refresh    chan struct{}
	lastStatus *BatteryStatus
}

func NewMonitor(cfg Config, onUpdate func(BatteryStatus)) *Monitor {
	return &Monitor{
		cfg:      cfg,
		onUpdate: onUpdate,
		stop:     make(chan struct{}),
		refresh:  make(chan struct{}, 1),
	}
}

func (m *Monitor) Start() {
	go m.loop()
}

func (m *Monitor) Stop() {
	close(m.stop)
}

// Refresh は次のポーリング周期を待たずに即時更新を要求する
func (m *Monitor) Refresh() {
	select {
	case m.refresh <- struct{}{}:
	default:
	}
}

func (m *Monitor) loop() {
	m.poll() // 起動直後に即時ポーリング

	for {
		interval := m.nextInterval()
		timer := time.NewTimer(time.Duration(interval) * time.Second)

		select {
		case <-m.stop:
			timer.Stop()
			return
		case <-m.refresh:
			timer.Stop()
			m.poll()
		case <-timer.C:
			m.poll()
		}
	}
}

func (m *Monitor) poll() {
	if m.cfg.FPSGuard.Enabled && isGameActive(m.cfg.FPSGuard.GameProcesses) {
		log.Println("ゲームがアクティブなためポーリングをスキップ")
		return
	}

	dev := findRazerMouse()
	if dev == nil {
		return
	}

	status, err := readBattery(dev.Path, dev.TransactionID, m.cfg.Debug.DumpRawBytes)
	if err != nil {
		log.Printf("バッテリー取得エラー (%s): %v", dev.Name, err)
		return
	}

	m.lastStatus = status
	chargeLabel := ""
	if status.IsCharging {
		chargeLabel = " (充電中)"
	}
	log.Printf("%s — %d%%%s", dev.Name, status.Level, chargeLabel)
	m.onUpdate(*status)
}

func (m *Monitor) nextInterval() int {
	if m.lastStatus == nil {
		return m.cfg.Monitor.PollIntervalNormal
	}
	if m.lastStatus.IsCharging {
		return m.cfg.Monitor.PollIntervalCharging
	}
	if m.lastStatus.Level <= m.cfg.Alert.LowBatteryThreshold {
		return m.cfg.Monitor.PollIntervalLow
	}
	return m.cfg.Monitor.PollIntervalNormal
}
