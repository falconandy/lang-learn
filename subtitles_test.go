package langlearn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLMarkupCleaner(t *testing.T) {
	c := NewHTMLMarkupCleaner()

	assert.Equal(t, "", c.Clean(""))
	assert.Equal(t, "test", c.Clean("test"))
	assert.Equal(t, "test", c.Clean("<b>test"))
	assert.Equal(t, "test", c.Clean("test</b>"))
	assert.Equal(t, "test", c.Clean("te{i}st"))
	assert.Equal(t, "test", c.Clean("tes{/i}t"))
	assert.Equal(t, "test", c.Clean("<font color='red'>te</font>st"))
}
