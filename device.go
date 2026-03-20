package main

import (
	"log"

	hid "github.com/sstallion/go-hid"
)

// RazerDevice は検出されたRazerデバイスの情報を保持する
type RazerDevice struct {
	Name          string
	PID           uint16
	Path          string
	TransactionID byte
}

// findRazerMouse は接続中のRazerワイヤレスマウスを返す
func findRazerMouse() *RazerDevice {
	var found *RazerDevice

	hid.Enumerate(razerVID, 0, func(info *hid.DeviceInfo) error {
		// コントロールインターフェース (interface_number == 0) のみ対象
		if info.InterfaceNbr != 0 {
			return nil
		}
		if found != nil {
			return nil // 最初に見つかったデバイスのみ使用
		}

		entry, ok := knownMice[info.ProductID]
		if !ok {
			return nil
		}

		found = &RazerDevice{
			Name:          entry.name,
			PID:           info.ProductID,
			Path:          info.Path,
			TransactionID: entry.transactionID,
		}
		log.Printf("デバイス検出: %s (PID=0x%04X)", found.Name, found.PID)
		return nil
	})

	if found == nil {
		log.Println("Razerワイヤレスマウスが見つかりませんでした")
	}
	return found
}
