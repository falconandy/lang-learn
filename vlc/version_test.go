package vlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionGet(t *testing.T) {
	factory := NewVersionFactory()

	assert.Equal(t, "4.0.0", factory.Get("5.0.0").version.String())
	assert.Equal(t, "4.0.0", factory.Get("4.0.0").version.String())
	assert.Equal(t, "0.0.0", factory.Get("3.9.0").version.String())
	assert.Equal(t, "0.0.0", factory.Get("0.0.0").version.String())
}

func TestVersionFind(t *testing.T) {
	factory := NewVersionFactory()

	text := `Command Line Interface initialized. Type 'help' for help.
VLC media player 4.1.3
Command Line Interface initialized. Type 'help' for help.
`

	assert.Equal(t, "4.0.0", factory.Find(text).version.String())
}
