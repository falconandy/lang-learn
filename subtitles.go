package langlearn

import (
	"regexp"
	"time"
)

type Subtitles []*Subtitle

type Subtitle struct {
	Start time.Duration
	End   time.Duration
	Text  []string
}

type htmlMarkupCleaner struct {
	tagRe *regexp.Regexp
}

func newHTMLMarkupCleaner() *htmlMarkupCleaner {
	return &htmlMarkupCleaner{
		tagRe: regexp.MustCompile(`[<{]\s*/?\s*(b|u|i|font)\b[^>}]*[>}]`),
	}
}

func (c *htmlMarkupCleaner) Clean(subtitleLine string) string {
	return c.tagRe.ReplaceAllLiteralString(subtitleLine, "")
}
