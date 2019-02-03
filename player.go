package langlearn

type Position int

type VideoPlayer interface {
	Start() (progress <-chan Position, err error)
	Shutdown() error

	Play(videoPath string) error
	Seek(position Position) error

	SpeedSlower()
	SpeedNormal()
}
