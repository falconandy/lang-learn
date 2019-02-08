package srt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/falconandy/lang-learn"
)

const (
	timePattern = `((\d{0,2}):(\d{0,2}):(\d{0,2}),(\d{0,3}))`
)

type Parser struct {
	utf8Signature []byte
	indexRe       *regexp.Regexp
	timeRe        *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		utf8Signature: []byte{'\xEF', '\xBB', '\xBF'},
		indexRe:       regexp.MustCompile(`^\d+$`),
		timeRe:        regexp.MustCompile(fmt.Sprintf(`^\s*%s --> %s\s*$`, timePattern, timePattern)),
	}
}

func (p *Parser) Parse(r io.Reader, cleaner langlearn.SubtitleCleaner) (subtitles []*langlearn.Subtitle, err error) {
	s := bufio.NewScanner(r)
	var subtitle *langlearn.Subtitle
	firstLine := true
	for s.Scan() {
		line := s.Bytes()

		if firstLine {
			line = bytes.TrimPrefix(line, p.utf8Signature)
			firstLine = false
		}

		switch {
		case len(line) == 0:
			subtitle = nil
		case subtitle == nil:
			if p.indexRe.Match(line) {
				subtitle = &langlearn.Subtitle{}
				subtitles = append(subtitles, subtitle)
			} else {
				return nil, fmt.Errorf("incorrect index format: %s", string(line))
			}
		default:
			if subtitle.End > 0 {
				line := string(line)
				if cleaner != nil {
					line = cleaner.Clean(line)
				}
				subtitle.Text = append(subtitle.Text, line)
			} else {
				match := p.timeRe.FindSubmatch(line)
				if match == nil {
					return nil, fmt.Errorf("incorrect time format: %s", string(line))
				}
				subtitle.Start = p.parseDuration(match[2:6])
				subtitle.End = p.parseDuration(match[7:11])
			}
		}
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return subtitles, nil
}

func (p *Parser) parseDuration(parts [][]byte) time.Duration {
	h, _ := strconv.Atoi(string(parts[0]))
	m, _ := strconv.Atoi(string(parts[1]))
	s, _ := strconv.Atoi(string(parts[2]))
	ms, _ := strconv.Atoi(string(parts[3]))
	return time.Hour*time.Duration(h) + time.Minute*time.Duration(m) +
		time.Second*time.Duration(s) + time.Millisecond*time.Duration(ms)
}
