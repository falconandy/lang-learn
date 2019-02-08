package langlearn

import (
	"regexp"
	"time"
)

type Subtitle struct {
	Start time.Duration
	End   time.Duration
	Text  []string
}

type SubtitleCleaner interface {
	Clean(subtitleLine string) string
}

type SubtitleIndex interface {
	BuildIndex(subtitles []*Subtitle)
	Find(word string) []*Subtitle
}

type htmlMarkupCleaner struct {
	tagRe *regexp.Regexp
}

func NewHTMLMarkupCleaner() SubtitleCleaner {
	return &htmlMarkupCleaner{
		tagRe: regexp.MustCompile(`[<{]\s*/?\s*(b|u|i|font)\b[^>}]*[>}]`),
	}
}

func (c *htmlMarkupCleaner) Clean(subtitleLine string) string {
	return c.tagRe.ReplaceAllLiteralString(subtitleLine, "")
}
