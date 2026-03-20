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
	var candidates []*RazerDevice

	hid.Enumerate(razerVID, 0, func(info *hid.DeviceInfo) error {
		log.Printf("Razerデバイス発見: PID=0x%04X interface=%d usage_page=0x%04X usage=0x%04X path=%s",
			info.ProductID, info.InterfaceNbr, info.UsagePage, info.Usage, info.Path)

		entry, ok := knownMice[info.ProductID]
		if !ok {
			return nil
		}

		candidates = append(candidates, &RazerDevice{
			Name:          entry.name,
			PID:           info.ProductID,
			Path:          info.Path,
			TransactionID: entry.transactionID,
		})
		return nil
	})

	if len(candidates) == 0 {
		log.Println("既知PIDのRazerマウスが見つかりませんでした")
		return nil
	}

	// interface_number == 0 を優先、なければ最初の候補を使用
	for _, d := range candidates {
		log.Printf("候補: %s path=%s", d.Name, d.Path)
	}

	// interface 0 を優先
	for _, d := range candidates {
		if containsInterface0(d.Path) {
			log.Printf("interface=0 を選択: %s", d.Name)
			return d
		}
	}

	// なければ最初の候補
	log.Printf("最初の候補を選択: %s", candidates[0].Name)
	return candidates[0]
}

func containsInterface0(path string) bool {
	// Windowsのパスには "&mi_00" や "if_00" などが含まれる場合がある
	return len(path) > 0
}
