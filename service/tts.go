package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

var (
	dll            *syscall.DLL
	initTTSProc    *syscall.Proc
	speakTextProc  *syscall.Proc
	releaseTTSProc *syscall.Proc
)

func init() {
	var err error
	dll, err = syscall.LoadDLL("tts.dll")
	if err != nil {
		panic(fmt.Sprintf("无法加载 DLL: %v", err))
	}
	initTTSProc, _ = dll.FindProc("initTTS")
	speakTextProc, _ = dll.FindProc("speakText")
	releaseTTSProc, _ = dll.FindProc("releaseTTS")
}

func initTTS() {
	initTTSProc.Call()
}

func speakText(rate float64, text string) {
	text = `<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='zh-CN'>
<prosody rate='` + strconv.FormatFloat(rate, 'f', -1, 64) + `'>` + text + `</prosody>
</speak>`

	textPtr, _ := syscall.UTF16PtrFromString(text)
	speakTextProc.Call(uintptr(unsafe.Pointer(textPtr)))
}

func releaseTTS() {
	releaseTTSProc.Call()
}
