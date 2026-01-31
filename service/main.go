package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	port int = 13520

	subtitles     []Subtitle = make([]Subtitle, 0)
	channel       chan int64 = make(chan int64, 1024)
	subtitleIndex int        = 0

	lastTime   int64 = 0
	isSpeaking bool  = false
	adjustTime bool  = false
)

func main() {
	fmt.Println("启动中...")
	initTTS() // 初始化TTS
	time.Sleep(3 * time.Second)
	speakText(1.0, "启动成功")
	fmt.Println("启动成功...")

	// 启动服务器
	go StartServer()

	// 监听键盘输入
	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				close(channel)
				return
			}
			if buf[0] == 'q' {
				close(channel)
				return
			}
		}
	}()

	// 播放字幕
	for time := range channel {
		x := time - lastTime
		lastTime = time
		fmt.Printf("\rtime: %s", timeToString(time))

		// 如果正在调整时间，则跳过
		if adjustTime {
			continue
		}

		if x < -1000 {
			adjustTime = true // 暂停播放
			fmt.Printf("\n倒退: %s\n", timeToString(x))
			for subtitleIndex < len(subtitles) && subtitleIndex > 0 {
				//处理倒退的情况
				subtitle := subtitles[subtitleIndex]
				if subtitle.StartMs > time { // 如果字幕开始时间大于当前时间，则倒退
					subtitleIndex--
				} else {
					break
				}
			}
			adjustTime = false // 恢复播放
			continue
		}

		if x > 1000 {
			adjustTime = true // 暂停播放
			fmt.Printf("\n跳过: %s\n", timeToString(x))
			for subtitleIndex < len(subtitles) {
				// 处理跳过的情况
				subtitle := subtitles[subtitleIndex]
				if subtitle.StartMs < time { // 如果字幕开始时间小于当前时间，则跳过
					subtitleIndex++
				} else {
					break
				}
			}
			adjustTime = false // 恢复播放
			continue
		}

		// 如果字幕索引大于字幕长度，则跳过 (处理视频播放完毕)
		if subtitleIndex >= len(subtitles) {
			continue
		}

		// 如果正在播放字幕，则跳过 (处理插件无脑发送时间)
		if isSpeaking {
			continue
		}
		isSpeaking = true

		go func() {
			text := ""
			speakRate := 1.0
			for {
				// 拼接多个字幕
				subtitle := subtitles[subtitleIndex]
				if subtitle.StartMs < time {
					text += subtitle.Text
					subtitleIndex++

					// 如果语速小于5，则增加语速
					speakRate += 1
				} else {
					break
				}
			}

			if text != "" {
				speakText(speakRate, text) // 播放字幕
			}

			isSpeaking = false
		}()
	}
	releaseTTS() // 释放TTS
	fmt.Println("\n程序结束")
}

func StartServer() {
	http.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		var youtubeSubtitle YoutubeSubtitle
		err = json.Unmarshal(body, &youtubeSubtitle)
		if err != nil {
			return
		}

		subtitles = youtubeSubtitle.ToSubtitle()
		lastTime = 0
		subtitleIndex = 0
		fmt.Println("\n加载字幕")
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		time := r.URL.Query().Get("time")

		doubleTime, err := strconv.ParseFloat(time, 64)
		if err != nil {
			return
		}

		channel <- int64(doubleTime * 1000)
	})

	fmt.Printf("启动服务器: http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func timeToString(time int64) string {
	hours := time / 3600000
	minutes := (time % 3600000) / 60000
	seconds := (time % 60000) / 1000
	milliseconds := time % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
