package vlc

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/falconandy/lang-learn"
)

const (
	startupDelay       = time.Second * 2
	initialReadTimeout = time.Millisecond * 500
	nextReadTimeout    = time.Millisecond * 200
)

type vlcPlayer struct {
	exePath  string
	tcpPort  int
	promptRe *regexp.Regexp
	version  *Version

	conn       net.Conn
	connReader *bufio.Reader
	commands   chan<- *command
}

func NewPlayer(exePath string, tcpPort int) langlearn.VideoPlayer {
	return &vlcPlayer{
		exePath:  exePath,
		tcpPort:  tcpPort,
		promptRe: regexp.MustCompile(`(>\s+)+`),
	}
}

func (p *vlcPlayer) Start() (<-chan langlearn.Position, error) {
	cmd := exec.Command(p.exePath,
		"--extraintf=rc",
		fmt.Sprintf("--rc-host=%s:%d", "localhost", p.tcpPort),
		"--one-instance") //"--rc-quiet",
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	time.Sleep(startupDelay)

	p.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", p.tcpPort))
	if err != nil {
		return nil, err
	}
	p.connReader = bufio.NewReader(p.conn)

	helpOutput, err := p.execCommand("help")
	if err != nil {
		return nil, err
	}

	p.version = NewVersionFactory().Find(helpOutput)

	commands, progress := make(chan *command), make(chan langlearn.Position)
	p.commands = commands
	go p.run(commands, progress)

	return progress, nil
}

func (p *vlcPlayer) Shutdown() error {
	cmd := newCommand(p, p.version.shutdownCommand)
	p.commands <- cmd
	result := <-cmd.result
	return result.err
}

func (p *vlcPlayer) Play(videoPath string) error {
	addCmd := newCommand(p, fmt.Sprintf(`add %s`, videoPath))
	p.commands <- addCmd
	addResult := <-addCmd.result
	if addResult.err != nil {
		return addResult.err
	}

	strackCmd := newCommand(p, "strack -1")
	p.commands <- strackCmd
	strackResult := <-strackCmd.result
	return strackResult.err
}

func (p *vlcPlayer) Seek(position langlearn.Position) error {
	cmd := newCommand(p, fmt.Sprintf("seek %d", position))
	p.commands <- cmd
	result := <-cmd.result
	return result.err
}

func (*vlcPlayer) SpeedSlower() {
	panic("implement me")
}

func (*vlcPlayer) SpeedNormal() {
	panic("implement me")
}

func (p *vlcPlayer) run(commands <-chan *command, progress chan<- langlearn.Position) {
	defer close(progress)
	defer func() { _ = p.conn.Close() }()

	currentPosition := -1

LOOP:
	for {
		select {
		case cmd := <-commands:
			cmd.Execute()
			if cmd.cmd == p.version.shutdownCommand {
				break LOOP
			}
		case <-time.After(time.Millisecond * 100):
			positionResponse, err := p.execCommand("get_time")
			if err != nil {
				fmt.Printf("failed to execute a command '%s': %v\n", "get_time", err)
				return
			}

			if positionResponse != "" {
				position, err := strconv.Atoi(positionResponse)
				if err != nil {
					fmt.Printf("position '%s' isn't a number: %v\n", positionResponse, err)
					continue
				}

				if position != currentPosition {
					select {
					case progress <- langlearn.Position(position):
						currentPosition = position
					default:
					}
				}
			}
		}
	}
}

func (p *vlcPlayer) execCommand(command string) (string, error) {
	println("INP: ", command)

	_, err := fmt.Fprintln(p.conn, command)
	if err != nil {
		return "", err
	}

	var output []string
	readTimeout := initialReadTimeout
	for {
		err := p.conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			fmt.Printf("can't set read deadline for a VLC connection: %v\n", err)
			break
		}

		line, err := p.connReader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(p.promptRe.ReplaceAllLiteralString(line, ""))

		println(line)

		output = append(output, line)
		readTimeout = nextReadTimeout
	}
	return strings.Join(output, "\n"), nil
}