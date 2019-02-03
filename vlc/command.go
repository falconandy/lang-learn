package vlc

type commandResult struct {
	output string
	err    error
}

type command struct {
	player *vlcPlayer
	cmd    string
	result chan *commandResult
}

func newCommand(player *vlcPlayer, cmd string) *command {
	return &command{
		player: player,
		cmd:    cmd,
		result: make(chan *commandResult, 1),
	}
}

func (c *command) Execute() {
	output, err := c.player.execCommand(c.cmd)
	c.result <- &commandResult{output: output, err: err}
	close(c.result)
}

func (c *command) String() string {
	return c.cmd
}
