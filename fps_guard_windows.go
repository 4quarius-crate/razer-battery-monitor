//go:build windows

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

// setLowPriority はプロセス優先度を BELOW_NORMAL に設定する
func setLowPriority() {
	handle, err := windows.GetCurrentProcess()
	if err != nil {
		log.Printf("プロセスハンドル取得失敗: %v", err)
		return
	}
	const belowNormal = 0x00004000
	r, _, err := kernel32.NewProc("SetPriorityClass").Call(uintptr(handle), belowNormal)
	if r == 0 {
		log.Printf("プロセス優先度の設定失敗: %v", err)
		return
	}
	log.Printf("プロセス優先度を BELOW_NORMAL に設定しました (PID=%d)", os.Getpid())
}

// foregroundProcessName はフォアグラウンドウィンドウのプロセス名（小文字）を返す
func foregroundProcessName() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return ""
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))
	r, _, _ := procQueryFullProcessImageName.Call(
		uintptr(handle),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if r == 0 {
		return ""
	}

	fullPath := windows.UTF16ToString(buf[:size])
	return strings.ToLower(filepath.Base(fullPath))
}

// isGameActive はゲームプロセスがフォアグラウンドにあるか判定する
func isGameActive(gameProcesses []string) bool {
	if len(gameProcesses) == 0 {
		return false
	}
	name := foregroundProcessName()
	if name == "" {
		return false
	}
	for _, g := range gameProcesses {
		if strings.ToLower(g) == name {
			return true
		}
	}
	return false
}
