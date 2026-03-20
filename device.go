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
	InterfaceNbr  int
}

// findRazerMouse は接続中のRazerワイヤレスマウスを返す
func findRazerMouse() *RazerDevice {
	var candidates []*RazerDevice

	hid.Enumerate(razerVID, 0, func(info *hid.DeviceInfo) error {
		log.Printf("HIDデバイス: PID=0x%04X interface=%d usage_page=0x%04X path=%s",
			info.ProductID, info.InterfaceNbr, info.UsagePage, info.Path)

		entry, ok := knownMice[info.ProductID]
		if !ok {
			return nil
		}

		candidates = append(candidates, &RazerDevice{
			Name:          entry.name,
			PID:           info.ProductID,
			Path:          info.Path,
			TransactionID: entry.transactionID,
			InterfaceNbr:  info.InterfaceNbr,
		})
		return nil
	})

	if len(candidates) == 0 {
		log.Println("既知PIDのRazerマウスが見つかりませんでした")
		return nil
	}

	for _, d := range candidates {
		log.Printf("候補: %s interface=%d path=%s", d.Name, d.InterfaceNbr, d.Path)
	}

	// usage_page=0xFF00（ベンダー固有）のインターフェースを優先
	// なければ interface 2 → 0 の順で試す
	for _, preferredIf := range []int{2, 1, 0} {
		for _, d := range candidates {
			if d.InterfaceNbr == preferredIf {
				log.Printf("interface=%d を選択: %s", d.InterfaceNbr, d.Name)
				return d
			}
		}
	}

	log.Printf("最初の候補を選択: %s interface=%d", candidates[0].Name, candidates[0].InterfaceNbr)
	return candidates[0]
}
