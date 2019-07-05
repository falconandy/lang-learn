package bleve

import (
	"strings"

	"github.com/blevesearch/segment"

	langlearn "github.com/falconandy/lang-learn"
)

type subtitleIndex struct {
	index map[string][]*langlearn.Subtitle
}

func NewSubtitleIndex() langlearn.SubtitleIndex {
	return &subtitleIndex{}
}

func (idx *subtitleIndex) BuildIndex(subtitles []*langlearn.Subtitle) {
	index := make(map[string][]*langlearn.Subtitle)

	for _, s := range subtitles {
		words := make(map[string]bool)
		for _, t := range s.Text {
			seg := segment.NewWordSegmenter(strings.NewReader(t))
			for seg.Segment() {
				if seg.Type() != segment.Letter {
					continue
				}
				word := strings.ToLower(seg.Text())
				if strings.HasSuffix(word, "'s") && word != "it's" && word != "he's" && word != "she's" {
					word = strings.TrimSuffix(word, "'s")
				}

				if strings.ContainsAny(word, ".,") {
					continue
				}

				words[word] = true
			}
		}
		for word := range words {
			index[word] = append(index[word], s)
		}
	}

	idx.index = index
}

func (idx *subtitleIndex) Find(word string) []*langlearn.Subtitle {
	return idx.index[strings.ToLower(word)]
}
