package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"

	langlearn "github.com/falconandy/lang-learn"
)

func TestSubtitleIndex(t *testing.T) {
	subtitles := []*langlearn.Subtitle{
		{
			Text: []string{
				"it's a first line",
				"it's a second line",
			},
		},
		{
			Text: []string{
				"first line",
				"parent's house",
			},
		},
		{
			Text: []string{
				"row www.google.com",
				"second row",
			},
		},
	}

	idx := NewSubtitleIndex()
	idx.BuildIndex(subtitles)
	result := idx.Find("second")

	assert.Equal(t, 2, len(result))
	assert.Equal(t, subtitles[0], result[0])
	assert.Equal(t, subtitles[2], result[1])

	result = idx.Find("google")
	assert.Equal(t, 0, len(result))

	result = idx.Find("www.google.com")
	assert.Equal(t, 0, len(result))

	result = idx.Find("parent's")
	assert.Equal(t, 0, len(result))

	result = idx.Find("parent")
	assert.Equal(t, 1, len(result))
	assert.Equal(t, subtitles[1], result[0])
}
