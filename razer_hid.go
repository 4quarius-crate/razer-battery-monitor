package main

import (
	"fmt"
	"log"

	hid "github.com/sstallion/go-hid"
)

const (
	razerVID  = 0x1532
	reportLen = 91
)

// knownMice は対応しているRazerワイヤレスマウスのPIDリスト
var knownMice = map[uint16]mouseEntry{
	0x0177: {"Razer Basilisk Ultimate (Wired/Receiver)", 0x1F},
	0x0178: {"Razer Basilisk Ultimate (Wireless)", 0x1F},
	0x007C: {"Razer DeathAdder V2 Pro (Wired)", 0x1F},
	0x007D: {"Razer DeathAdder V2 Pro (Wireless)", 0x1F},
	0x0152: {"Razer Viper Ultimate (Wired)", 0x1F},
	0x0153: {"Razer Viper Ultimate (Wireless)", 0x1F},
	0x0192: {"Razer Naga Pro (Wired)", 0x1F},
	0x0193: {"Razer Naga Pro (Wireless)", 0x1F},
}

type mouseEntry struct {
	name          string
	transactionID byte
}

// BatteryStatus はバッテリー状態を保持する
type BatteryStatus struct {
	Level      int
	IsCharging bool
}

// buildReport は90バイトのRazer HIDフィーチャーレポートを組み立てる
func buildReport(cmdClass, cmdID, transactionID byte) []byte {
	buf := make([]byte, reportLen)
	buf[0] = 0x00 // Report ID
	buf[1] = 0x00 // Status
	buf[2] = transactionID
	buf[3] = 0x00 // Remaining packets Hi
	buf[4] = 0x00 // Remaining packets Lo
	buf[5] = 0x00 // Protocol type
	buf[6] = 0x02 // Data size
	buf[7] = cmdClass
	buf[8] = cmdID
	// CRC: buf[2]〜buf[87] の XOR
	var crc byte
	for _, b := range buf[2:88] {
		crc ^= b
	}
	buf[89] = crc
	return buf
}

// readBattery はデバイスを開いてバッテリー情報を取得する
func readBattery(path string, transactionID byte, dumpRaw bool) (*BatteryStatus, error) {
	dev, err := hid.OpenPath(path)
	if err != nil {
		return nil, fmt.Errorf("HIDデバイスを開けませんでした: %w", err)
	}
	defer dev.Close()

	// バッテリーレベル取得
	report := buildReport(0x02, 0x07, transactionID)
	if dumpRaw {
		log.Printf("TX (battery level): %X", report)
	}
	if _, err := dev.SendFeatureReport(report); err != nil {
		return nil, fmt.Errorf("フィーチャーレポート送信エラー: %w", err)
	}
	resp := make([]byte, reportLen)
	resp[0] = 0x00
	if _, err := dev.GetFeatureReport(resp); err != nil {
		return nil, fmt.Errorf("フィーチャーレポート受信エラー: %w", err)
	}
	if dumpRaw {
		log.Printf("RX (battery level): %X", resp)
	}
	level := parseBatteryLevel(resp)

	// 充電状態取得
	chargeReport := buildReport(0x02, 0x09, transactionID)
	dev.SendFeatureReport(chargeReport)
	chargeResp := make([]byte, reportLen)
	chargeResp[0] = 0x00
	dev.GetFeatureReport(chargeResp)
	charging := chargeResp[9]&0x01 != 0

	return &BatteryStatus{Level: level, IsCharging: charging}, nil
}

func parseBatteryLevel(resp []byte) int {
	if len(resp) < 11 {
		return -1
	}
	raw := resp[9]
	if raw == 0 || raw == 0xFF {
		raw = resp[10] // 一部デバイスはbuf[10]に値を格納する
	}
	level := int(float64(raw) / 255.0 * 100.0)
	if level > 100 {
		level = 100
	}
	return level
}
