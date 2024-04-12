package cli

import "os"

type FlagType uint

type CLI struct {
	Commands []*Command
	Help     string
}

type Command struct {
	Name    string
	Help    string
	Flags   []*Flag
	Handler func()
}

type Flag struct {
	Name string
	Type *FlagType
	Help string
}

const (
	FlagInt    = FlagType(0)
	FlagString = FlagType(1)
	FlagBool   = FlagType(2)
)

func (c *CLI) Handle() {
	if len(os.Args) == 1 {
		c.basicHelp()
		return
	}
	id := os.Args[1]
	f := c.findCommand(id)
	if f != nil {
		f.Handler()
		return
	}
	if id == "help" && len(os.Args) == 2 {
		c.basicHelp()
	} else if id == "help" {
		sub := os.Args[2]
		f = c.findCommand(sub)
		if f != nil {
			c.commandHelp(f)
			return
		}
		c.invalidCommand()
	} else {
		c.invalidCommand()
	}
}

func (c *CLI) findCommand(name string) *Command {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

func (c *CLI) basicHelp() {

}

func (c *CLI) commandHelp(cmd *Command) {

}

func (c *CLI) invalidCommand() {

}
