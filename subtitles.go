package langlearn

import (
	"time"
)

type Subtitles []*Subtitle

type Subtitle struct {
	Start time.Duration
	End   time.Duration
	Text  []string
}
