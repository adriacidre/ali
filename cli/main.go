package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adriacidre/aliases/config"
	"github.com/fatih/color"
	. "github.com/logrusorgru/aurora"
	"github.com/zyedidia/highlight"
)

func main() {
	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("homedir: %s", err)
	}

	path := home + "/.aliases"
	c := config.Config{}
	c.Parse(path)

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "list", "l":
		list(c)
	case "details", "d":
		details(c, os.Args[2])
	case "add", "a":
		add(c, path, os.Args[2])
	case "rm", "r":
		rm(c, path, os.Args[2])
	case "source", "s":
		source(c, path)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func list(c config.Config) {
	groups := make(map[string][]string)

	for _, c := range c.Commands {
		cmd := fmt.Sprint("  ", Yellow(c.Name), " ["+c.Usage+"]:", c.Description)
		group := "default"
		if c.Group != "" && c.Group != "<group>" {
			group = c.Group
		}

		if _, ok := groups[group]; !ok {
			groups[group] = make([]string, 0)
		}

		groups[group] = append(groups[group], cmd)

	}
	for k, v := range groups {
		fmt.Println(Magenta(k))
		for _, val := range v {
			fmt.Println(val)
		}
		fmt.Println("")
	}
}

func details(c config.Config, name string) {
	for _, c := range c.Commands {
		if c.Name == name {
			fmt.Println(Yellow(c.Name), ":", c.Description)
			fmt.Println("  usage: " + c.Usage)
			if len(c.ID) > 0 {
				fmt.Println("  id: " + c.ID)
			}
			if len(c.Tags) > 0 {
				fmt.Println("  tags: " + strings.Join(c.Tags, ","))
			}
			if len(c.Depends) > 0 {
				fmt.Println("  dependencies: " + strings.Join(c.Depends, ","))
			}

			if len(c.Content) > 0 {

				printHighlighted(c.Content)

			}
		}
	}
}

func add(c config.Config, path, name string) {
	alias, err := CaptureInputFromEditor(
		GetPreferredEditorFromEnvironment,
		c.HeaderTemplate(name),
	)
	if err != nil {
		panic(err)
	}

	ctmp := config.Config{}
	ctmp.ParseBody(alias)

	c.AddCommand(ctmp.Commands[0])
	c.Save(path)

}

func rm(c config.Config, path, name string) {
	c.RmCommand(name)
	c.Save(path)
}

func source(c config.Config, path string) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	fmt.Println(text)
}

func printHighlighted(input string) {
	syntaxFile, _ := ioutil.ReadFile("sh.yaml")
	syntaxDef, err := highlight.ParseDef(syntaxFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Make a new highlighter from the definition
	h := highlight.NewHighlighter(syntaxDef)
	matches := h.HighlightString(input)

	lines := strings.Split(input, "\n")
	for lineN, l := range lines {
		for colN, c := range l {
			// Check if the group changed at the current position
			if group, ok := matches[lineN][colN]; ok {
				// Check the group name and set the color accordingly (the colors chosen are arbitrary)
				if group == highlight.Groups["statement"] {
					color.Set(color.FgGreen)
				} else if group == highlight.Groups["preproc"] {
					color.Set(color.FgHiRed)
				} else if group == highlight.Groups["special"] {
					color.Set(color.FgBlue)
				} else if group == highlight.Groups["constant.string"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["constant.specialChar"] {
					color.Set(color.FgHiMagenta)
				} else if group == highlight.Groups["type"] {
					color.Set(color.FgYellow)
				} else if group == highlight.Groups["constant.number"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["comment"] {
					color.Set(color.FgHiGreen)
				} else {
					color.Unset()
				}
			}
			// Print the character
			fmt.Print(string(c))
		}
		// This is at a newline, but highlighting might have been turned off at the very end of the line so we should check that.
		if group, ok := matches[lineN][len(l)]; ok {
			if group == highlight.Groups["default"] || group == highlight.Groups[""] {
				color.Unset()
			}
		}

		fmt.Print("\n")
	}

	fmt.Sprintln("\n")
}