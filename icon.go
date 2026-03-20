package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const iconSize = 32

// makeBatteryIcon はバッテリー残量を描画したICOバイト列を返す
func makeBatteryIcon(level int, charging bool) []byte {
	img := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))

	// 背景（ダーク）
	bg := color.RGBA{30, 30, 30, 255}
	for y := 0; y < iconSize; y++ {
		for x := 0; x < iconSize; x++ {
			img.Set(x, y, bg)
		}
	}

	// バッテリー残量に応じた色
	var barColor color.RGBA
	switch {
	case charging:
		barColor = color.RGBA{80, 160, 255, 255}
	case level > 50:
		barColor = color.RGBA{80, 200, 80, 255}
	case level > 20:
		barColor = color.RGBA{220, 180, 0, 255}
	default:
		barColor = color.RGBA{220, 50, 50, 255}
	}

	// バッテリーバー（下から上に向かって伸びる）
	const padding = 3
	barX1, barX2 := padding, iconSize-padding
	barY2 := iconSize - padding
	barHeight := int(float64(iconSize-padding*2) * float64(level) / 100.0)
	barY1 := barY2 - barHeight

	for y := barY1; y < barY2; y++ {
		for x := barX1; x < barX2; x++ {
			img.Set(x, y, barColor)
		}
	}

	// パーセント数値を描画
	text := fmt.Sprintf("%d", level)
	if charging {
		text = "CHG"
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}
	bounds, _ := d.BoundString(text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	x := (iconSize - w) / 2
	y := iconSize/2 + 5
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d.DrawString(text)

	// PNG エンコード → ICO にラップ
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return pngToICO(buf.Bytes(), iconSize)
}

// pngToICO はPNGバイト列をICOコンテナにラップする（Windows Vista+対応）
func pngToICO(pngData []byte, size int) []byte {
	var buf bytes.Buffer

	sizeByte := byte(size)
	if size >= 256 {
		sizeByte = 0
	}

	// ICONDIR (6 bytes)
	binary.Write(&buf, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // type: ICO
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // image count

	// ICONDIRENTRY (16 bytes)
	buf.WriteByte(sizeByte)                                       // width
	buf.WriteByte(sizeByte)                                       // height
	buf.WriteByte(0)                                              // color count
	buf.WriteByte(0)                                              // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1))            // planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))           // bit count
	binary.Write(&buf, binary.LittleEndian, uint32(len(pngData))) // bytes in res
	binary.Write(&buf, binary.LittleEndian, uint32(22))           // image offset (6+16)

	buf.Write(pngData)
	return buf.Bytes()
}
