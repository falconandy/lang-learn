package vlc

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/falconandy/lang-learn"
)

const (
	startupDelay = time.Second * 2
)

type vlcPlayer struct {
	exePath string
	tcpPort int

	conn     *tcpConnection
	commands chan<- *command
}

func NewPlayer(exePath string, tcpPort int) langlearn.VideoPlayer {
	if exePath == "" {
		switch runtime.GOOS {
		case "windows":
			exePath = `C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`
		default:
			exePath = "vlc"
		}
	}

	return &vlcPlayer{
		exePath: exePath,
		tcpPort: tcpPort,
	}
}

func (p *vlcPlayer) Start() (<-chan langlearn.Position, error) {
	args := []string{
		"--extraintf=rc",
		fmt.Sprintf("--rc-host=%s:%d", "localhost", p.tcpPort),
		"--one-instance",
	}

	// TODO: specific to VLC version?
	if runtime.GOOS == "windows" {
		args = append(args, "--rc-quiet")
	}

	cmd := exec.Command(p.exePath, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	time.Sleep(startupDelay)

	p.conn = newTCPConnection(p.tcpPort)
	err = p.conn.Open()
	if err != nil {
		return nil, err
	}

	commands, progress := make(chan *command), make(chan langlearn.Position)
	p.commands = commands
	go p.conn.run(commands, progress)

	return progress, nil
}

func (p *vlcPlayer) Shutdown() error {
	_, err := p.execCommand(p.conn.version.shutdownCommand)
	return err
}

func (p *vlcPlayer) Play(videoPath string) error {
	_, err := p.execCommand("stop")
	if err != nil {
		return err
	}

	_, err = p.execCommand(fmt.Sprintf(`add %s`, videoPath))
	if err != nil {
		return err
	}

	for {
		output, err := p.execCommand("is_playing")
		if err != nil {
			return err
		}
		if output == "1" {
			break
		}
	}

	_, err = p.execCommand("strack -1")
	return err
}

func (p *vlcPlayer) Pause() error {
	_, err := p.execCommand("pause")
	return err
}

func (p *vlcPlayer) Stop() error {
	_, err := p.execCommand("stop")
	return err
}

func (p *vlcPlayer) Seek(position langlearn.Position) error {
	_, err := p.execCommand(fmt.Sprintf("seek %d", position))
	return err
}

func (p *vlcPlayer) SpeedSlower() error {
	_, err := p.execCommand("slower")
	return err
}

func (p *vlcPlayer) SpeedFaster() error {
	_, err := p.execCommand("faster")
	return err
}

func (p *vlcPlayer) SpeedNormal() error {
	_, err := p.execCommand("normal")
	return err
}

func (p *vlcPlayer) AudioTrack() (int, error) {
	output, err := p.execCommand("atrack")
	if err != nil {
		return -1, err
	}
	track, err := strconv.Atoi(output)
	if err != nil {
		return -1, fmt.Errorf("can't convert '%s' to a number: %v", output, err)
	}
	return track, nil
}

func (p *vlcPlayer) SetAudioTrack(track int) error {
	_, err := p.execCommand(fmt.Sprintf("atrack %d", track))
	return err
}

func (p *vlcPlayer) execCommand(cmd string) (output string, err error) {
	c := newCommand(p, cmd)
	p.commands <- c
	result := <-c.result
	return result.output, result.err
}
