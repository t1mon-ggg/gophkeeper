package cli

import (
	"fmt"
	"strings"

	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
)

func (c *CLI) executor(in string) {
	prefix := strings.Split(livePrefixState.livePrefix, "/")
	switch in {
	case "":
		if livePrefixState.livePrefix == "" {
			livePrefixState.livePrefix = ">>> "
		}
		livePrefixState.isEnable = true
	case "save":
		if livePrefixState.livePrefix == "" {
			livePrefixState.livePrefix = ">>> "
		}
		c.save()
	case "quit":
		c.save()
		livePrefixState.livePrefix = "quit> "
		c.wg.Done()
	case "..":
		if len(prefix) > 1 {
			newprefix := []string{}
			for i := 0; i < len(prefix)-1; i++ {
				newprefix = append(newprefix, prefix[i])
			}
			livePrefixState.livePrefix = strings.Join(newprefix, "/") + "> "
		}
		if len(prefix) == 1 {
			livePrefixState.isEnable = false
			livePrefixState.livePrefix = ">>> "
		}
	default:
		line := strings.Split(in, " ")
		if cmd, ok := helpers.FindCommand(in); ok {
			switch cmd {
			case "roster":
				if livePrefixState.livePrefix == "user> " {
					fmt.Println("todo something here")
				} else {
					fmt.Println("no such command")
				}
			case "revoke":
				if livePrefixState.livePrefix == "user> " {
					fmt.Println("todo something here")
				} else {
					fmt.Println("no such command")
				}
			case "confirm":
				if livePrefixState.livePrefix == "user> " {
					fmt.Println("todo something here")
				} else {
					fmt.Println("no such command")
				}
			case "list":
				if livePrefixState.livePrefix == "cmd> " {
					c.list()
				} else {
					fmt.Println("no such command")
				}
			case "get":
				if livePrefixState.livePrefix == "cmd> " {
					opts := line[2:]
					c.get(line[1], opts...)
				} else {
					fmt.Println("no such command")
				}
			case "insert":
				if livePrefixState.livePrefix == "cmd> " {
					c.insert(in)

				} else {
					fmt.Println("no such command")
				}
			case "delete":
				if livePrefixState.livePrefix == "cmd> " {
					if len(line) == 2 {
						c.delete(line[1])
					}
				} else {
					fmt.Println("no such command")
				}
			case "view":
				if livePrefixState.livePrefix == "config> " {
					c.view()
				} else {
					fmt.Println("no such command")
				}
			case "version":
				if livePrefixState.livePrefix == "config> " {
					c.view()
				} else {
					fmt.Println("no such command")
				}
			case "status":
				if livePrefixState.livePrefix == ">>> " || livePrefixState.livePrefix == "" {
					if c.config.Mode == "standalone" {
						c.log().Info(nil, "status not availible. standalone mode active")
						return
					}
					c.status()
					return
				} else {
					fmt.Println("no such command")
				}
			}
			return
		}
		var found bool
		newprefix := []string{}
		if len(prefix) == 1 {
			p := strings.Split(prefix[0], ">")
			if len(p) == 4 {
				p = []string{}
			}
			if len(p) > 1 {
				for i := 0; i <= len(prefix)-1; i++ {
					newprefix = append(newprefix, prefix[i])
				}
				for _, t := range suggests[livePrefixState.livePrefix] {
					if t.Text == in {
						found = true
					}
				}
				if found {
					livePrefixState.livePrefix = strings.TrimSuffix(strings.TrimSpace(strings.Join(newprefix, "/")), ">") + "/" + strings.TrimSpace(in) + "> "
				} else {
					fmt.Println("no such command")
				}
			} else {
				for i := 0; i <= len(prefix)-1; i++ {
					newprefix = append(newprefix, prefix[i])
				}
				for _, t := range suggests[livePrefixState.livePrefix] {
					if t.Text == in {
						found = true
					}
				}
				if found {
					livePrefixState.livePrefix = strings.TrimSpace(in) + "> "
				} else {
					fmt.Println("no such command")
				}
			}
		}
		if len(prefix) > 1 {
			for i := 0; i <= len(prefix)-1; i++ {
				newprefix = append(newprefix, prefix[i])
			}
			for _, t := range suggests[livePrefixState.livePrefix] {
				if t.Text == in {
					found = true
				}
			}
			if found {
				livePrefixState.livePrefix = strings.TrimSuffix(strings.TrimSpace(strings.Join(newprefix, "/")), ">") + "/" + strings.TrimSpace(in) + "> "
			} else {
				fmt.Println("no such command")
			}
		}
	}
	livePrefixState.isEnable = true
}
