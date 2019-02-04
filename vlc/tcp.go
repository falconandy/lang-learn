package vlc

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/falconandy/lang-learn"
)

const (
	initialReadTimeout = time.Millisecond * 500
	nextReadTimeout    = time.Millisecond * 200
)

type tcpConnection struct {
	port     int
	promptRe *regexp.Regexp

	conn       net.Conn
	connReader *bufio.Reader
	version    *Version
}

func newTCPConnection(port int) *tcpConnection {
	return &tcpConnection{
		port:     port,
		promptRe: regexp.MustCompile(`(>\s+)+`),
	}
}

func (c *tcpConnection) Open() error {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", c.port))
	if err != nil {
		return err
	}
	c.connReader = bufio.NewReader(c.conn)

	helpOutput, err := c.execCommand("help")
	if err != nil {
		return err
	}

	c.version = NewVersionFactory().Find(helpOutput)

	return nil
}

func (c *tcpConnection) run(commands <-chan *command, progress chan<- langlearn.Position) {
	defer close(progress)
	defer func() { _ = c.conn.Close() }()

	currentPosition := -1

LOOP:
	for {
		select {
		case cmd := <-commands:
			output, err := c.execCommand(cmd.cmd)
			cmd.result <- &commandResult{output: output, err: err}
			close(cmd.result)

			if cmd.cmd == c.version.shutdownCommand {
				break LOOP
			}
		case <-time.After(time.Millisecond * 100):
			positionResponse, err := c.execCommand("get_time")
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

func (c *tcpConnection) execCommand(command string) (string, error) {
	println("INP: ", command)

	_, err := fmt.Fprintln(c.conn, command)
	if err != nil {
		return "", err
	}

	var output []string
	readTimeout := initialReadTimeout
	for {
		err := c.conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			fmt.Printf("can't set read deadline for a VLC connection: %v\n", err)
			break
		}

		line, err := c.connReader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(c.promptRe.ReplaceAllLiteralString(line, ""))

		if strings.HasPrefix(line, "status change:") {
			continue
		}

		println(line)

		output = append(output, line)
		readTimeout = nextReadTimeout
	}
	return strings.Join(output, "\n"), nil
}
