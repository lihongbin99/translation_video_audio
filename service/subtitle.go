package main

import (
	"regexp"
	"strings"
)

var (
	replaceRegex = regexp.MustCompile(`\[[^\]]*\]`)
)

type Subtitle struct {
	StartMs    int64  `json:"start_ms"`
	DurationMs int64  `json:"duration_ms"`
	Text       string `json:"text"`
	IsDone     bool   `json:"is_done"`
}

type YoutubeSubtitle struct {
	Events []Event `json:"events"`
}

type Event struct {
	StartMs    int64     `json:"tStartMs"`
	DurationMs int64     `json:"dDurationMs"`
	Segments   []Segment `json:"segs"`
}

type Segment struct {
	Text string `json:"utf8"`
}

func (s *YoutubeSubtitle) ToSubtitle() []Subtitle {
	subtitles := make([]Subtitle, 0)
	for _, event := range s.Events {
		text := ""
		for _, segment := range event.Segments {
			text += segment.Text
		}

		text = strings.ReplaceAll(text, "\n", "")
		text = replaceRegex.ReplaceAllString(text, "")

		if strings.TrimSpace(text) == "" {
			continue
		}

		startMs := event.StartMs
		// 按照句号分割文本
		sentences := strings.Split(text, "。")
		for j, sentence := range sentences {
			// 计算这个句子占整个时间的比例
			durationMs := event.DurationMs * int64(len(sentence)) / int64(len(text))

			// 如果当前是第一个句子，并且前一个句子没有完成，则拼接到前一个句子的文本后面
			i := len(subtitles)
			if i > 0 && j == 0 && !subtitles[i-1].IsDone {
				subtitles[i-1].Text += sentence
				subtitles[i-1].DurationMs += durationMs
				startMs += durationMs
			} else {
				// 否则创建一个新的子标题
				subtitles = append(subtitles, Subtitle{
					StartMs:    startMs,
					DurationMs: durationMs,
					Text:       sentence,
					IsDone:     false,
				})
				startMs += durationMs
			}
		}
	}
	return subtitles
}
