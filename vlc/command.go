package vlc

import (
	"math"
)

type commandResult struct {
	output string
	err    error
}

type command struct {
	player            *vlcPlayer
	cmd               string
	responseLineCount int
	responsePostfix   string

	result chan *commandResult
}

func newCommand(player *vlcPlayer, cmd string) *command {
	return &command{
		player: player,
		cmd:    cmd,
		result: make(chan *commandResult, 1),
	}
}

func (c *command) withResponseLineCount(responseLineCount int) *command {
	c.responseLineCount = responseLineCount
	return c
}

func (c *command) withResponsePostfix(postfix string) *command {
	c.responsePostfix = postfix
	c.responseLineCount = math.MaxInt64
	return c
}

func (c *command) Execute() {
	output, err := c.player.execCommand(c.cmd, c.responseLineCount, c.responsePostfix)
	c.result <- &commandResult{output: output, err: err}
	close(c.result)
}

func (c *command) String() string {
	return c.cmd
}
