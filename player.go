package langlearn

type Position int

type VideoPlayer interface {
	Start() error
	Shutdown() error

	Play(videoPath string) error

	Seek(position Position) error

	SpeedFaster() error
	SpeedSlower() error
	SpeedNormal() error
}
