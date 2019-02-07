package vlc

import (
	"time"

	"github.com/falconandy/vlc"

	"github.com/falconandy/lang-learn"
)

const (
	checkIsPlayingInterval = time.Millisecond * 100
)

type player struct {
	vlcPlayer *vlc.Player
}

func NewPlayer(exePath string, tcpPort int) langlearn.VideoPlayer {
	return &player{
		vlcPlayer: vlc.NewPlayer(&vlc.PlayerConfig{ExePath: exePath, TCPPort: tcpPort}),
	}
}

func (p *player) Start() error {
	return p.vlcPlayer.Start()
}

func (p *player) Shutdown() error {
	return p.vlcPlayer.Shutdown()
}

func (p *player) Play(videoPath string) error {
	err := p.vlcPlayer.Stop()
	if err != nil {
		return err
	}

	err = p.vlcPlayer.Play(videoPath)
	if err != nil {
		return err
	}

	for {
		isPlaying, err := p.vlcPlayer.IsPlaying()
		if err != nil {
			return err
		}
		if isPlaying {
			break
		}

		time.Sleep(checkIsPlayingInterval)
	}

	return p.vlcPlayer.SetSubtitleTrack(-1)
}

func (p *player) Seek(position langlearn.Position) error {
	return p.vlcPlayer.Seek(vlc.Duration(position))
}

func (p *player) SpeedSlower() error {
	return p.vlcPlayer.SpeedSlower()
}

func (p *player) SpeedFaster() error {
	return p.vlcPlayer.SpeedFaster()
}

func (p *player) SpeedNormal() error {
	return p.vlcPlayer.SpeedNormal()
}
