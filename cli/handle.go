package cli

import (
	"fmt"
	"os"
)

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
	println("RSS Goes Social - CLI Help")
	println("==========================")
	m := 0
	for _, cmd := range c.Commands {
		if m < len(cmd.Name) {
			m = len(cmd.Name)
		}
	}
	for _, cmd := range c.Commands {
		sep := ""
		for range -len(cmd.Name) {
			sep += " "
		}
		fmt.Printf("\n%s%s - %s\n", cmd.Name, sep, cmd.Help)
	}
}

func (c *CLI) commandHelp(cmd *Command) {
	d := "RSS Goes Social - " + cmd.Name + " Help"
	sep := ""
	for range len(d) {
		sep += "="
	}
	println(d + "\n" + sep)

	m := 0
	for _, f := range cmd.Flags {
		if m < len(f.Name) {
			m = len(f.Name)
		}
	}
	for _, f := range cmd.Flags {
		sep = ""
		for range m - len(f.Name) {
			sep += " "
		}
		fmt.Printf("%s%s - %s\n", f.Name, sep, f.Help)
	}
}

func (c *CLI) invalidCommand() {
	println("Invalid command. Check the help with rss-goes-social help or with rss-goes-social help {command}")
}
