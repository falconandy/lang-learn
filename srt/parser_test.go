package srt

import (
	"strings"
	"testing"

	"github.com/falconandy/lang-learn"
	"github.com/stretchr/testify/assert"
)

func TestSRT_Empty(t *testing.T) {
	r := strings.NewReader("")
	subtitles, err := NewParser().Parse(r, nil)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(subtitles))
}

func TestSRT_Simple(t *testing.T) {
	r := strings.NewReader(`1
00:01:52,680 --> 00:01:55,630
<i>Born of cold and winter air</i>

2
00:01:55,640 --> 00:01:59,990
MAN: ♪ TripleDent gum
WOMAN: # Will make you smile ♪

3
00:02:00,000 --> 00:02:03,150
<b>This icy force both foul and fair</b>
`)
	subtitles, err := NewParser().Parse(r, langlearn.NewHTMLMarkupCleaner())

	assert.Nil(t, err)
	assert.Equal(t, 3, len(subtitles))

	assert.Equal(t, int64((1*60+52)*1000+680)*1000000, subtitles[0].Start.Nanoseconds())
	assert.Equal(t, int64((1*60+55)*1000+630)*1000000, subtitles[0].End.Nanoseconds())
	assert.Equal(t, 1, len(subtitles[0].Text))
	assert.Equal(t, "Born of cold and winter air", subtitles[0].Text[0])

	assert.Equal(t, int64((1*60+55)*1000+640)*1000000, subtitles[1].Start.Nanoseconds())
	assert.Equal(t, int64((1*60+59)*1000+990)*1000000, subtitles[1].End.Nanoseconds())
	assert.Equal(t, 2, len(subtitles[1].Text))
	assert.Equal(t, "MAN: ♪ TripleDent gum", subtitles[1].Text[0])
	assert.Equal(t, "WOMAN: # Will make you smile ♪", subtitles[1].Text[1])

	assert.Equal(t, int64((2*60+0)*1000+0)*1000000, subtitles[2].Start.Nanoseconds())
	assert.Equal(t, int64((2*60+3)*1000+150)*1000000, subtitles[2].End.Nanoseconds())
	assert.Equal(t, 1, len(subtitles[2].Text))
	assert.Equal(t, "This icy force both foul and fair", subtitles[2].Text[0])
}
