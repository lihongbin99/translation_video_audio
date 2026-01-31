package main

import "strings"

type Subtitle struct {
	StartMs    int64  `json:"start_ms"`
	DurationMs int64  `json:"duration_ms"`
	Text       string `json:"text"`
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
		text = strings.ReplaceAll(text, "[音乐]", "")

		if strings.TrimSpace(text) == "" {
			continue
		}

		subtitles = append(subtitles, Subtitle{
			StartMs:    event.StartMs,
			DurationMs: event.DurationMs,
			Text:       text,
		})
	}
	return subtitles
}
