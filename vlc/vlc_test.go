package vlc

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	langlearn "github.com/falconandy/lang-learn"
)

func TestVLCStart(t *testing.T) {
	player := NewPlayer("", 2019)
	err := player.Start()
	assert.Nil(t, err)
	time.Sleep(time.Second * 3)
	err = player.Shutdown()
	assert.Nil(t, err)
}

func TestVLCPlay(t *testing.T) {
	rand.Seed(time.Now().Unix())
	player := NewPlayer("", 2019)
	err := player.Start()
	assert.Nil(t, err)

	err = player.Play("/mnt/media/falconandy/_VIDEO/SLR.avi")
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	for i := 0; i < 5; i++ {
		err := player.Seek(langlearn.Position(rand.Intn(5500)))
		assert.Nil(t, err)

		time.Sleep(time.Second * 5)
	}
	err = player.Shutdown()
	assert.Nil(t, err)
}
