package langlearn

type Position int

type VideoPlayer interface {
	Start() (progress <-chan Position, err error)
	Shutdown() error

	Play(videoPath string) error
	Stop() error
	Pause() error

	Seek(position Position) error

	SpeedFaster() error
	SpeedSlower() error
	SpeedNormal() error

	AudioTrack() (int, error)
	SetAudioTrack(track int) error
}
