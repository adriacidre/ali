package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const COMMAND_PREFIX = "# "
const PROPERTY_PREFIX = "## "

type Command struct {
	ID          string
	Name        string
	Description string
	Content     string
	Group       string
	Usage       string
	Tags        []string
	Depends     []string
}

func (c *Command) Byte() []byte {
	body := strings.Join(c.ToLines(), "\n")
	return []byte(body)
}

func (c *Command) ToLines() []string {
	lines := []string{}
	lines = append(lines, "")
	lines = append(lines, COMMAND_PREFIX+c.Name)
	lines = append(lines, PROPERTY_PREFIX+"description: "+c.Description)
	if c.ID != "" {
		lines = append(lines, PROPERTY_PREFIX+"id: "+c.ID)
	}
	if c.Usage != "" {
		lines = append(lines, PROPERTY_PREFIX+"usage: "+c.Usage)
	}
	if len(c.Tags) > 0 {
		lines = append(lines, PROPERTY_PREFIX+"tags: "+strings.Join(c.Tags, ","))
	}
	if len(c.Depends) > 0 {
		lines = append(lines, PROPERTY_PREFIX+"depends: "+strings.Join(c.Depends, ","))
	}
	if len(c.Group) > 0 {
		lines = append(lines, PROPERTY_PREFIX+"group: "+c.Group)
	}
	lines = append(lines, c.Content)

	return lines
}

type Config struct {
	Commands []*Command
}

func (c *Config) Get(name string) (*Command, error) {
	for _, c := range c.Commands {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}

func (c *Config) Update(name string, input *Command) error {
	for pos, cmd := range c.Commands {
		if cmd.Name == name {
			c.Commands[pos] = input
			return nil
		}
	}
	return errors.New("not found")
}

func (c *Config) Parse(path string) error {
	lines, err := c.readLines(path)
	if err != nil {
		return err
	}

	c.parseLines(lines)

	return nil
}

func (c *Config) ParseBody(body []byte) {
	lines := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
	c.parseLines(lines)
}

func (c *Config) parseLines(lines []string) {
	var currentCommand *Command

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, COMMAND_PREFIX) { // This is an alias
			currentCommand = &Command{Name: strings.TrimPrefix(trimmedLine, COMMAND_PREFIX)}
			c.Commands = append(c.Commands, currentCommand)
		} else if strings.HasPrefix(trimmedLine, PROPERTY_PREFIX) { // This is an alias property
			parts := strings.Split(strings.TrimPrefix(trimmedLine, PROPERTY_PREFIX), ":")
			content := strings.TrimSpace(parts[1])
			switch parts[0] {
			case "id":
				currentCommand.ID = content
			case "group":
				currentCommand.Group = content
			case "description":
				currentCommand.Description = content
			case "usage":
				currentCommand.Usage = content
			case "tags":
				currentCommand.Tags = strings.Split(content, ",")
			case "depends":
				currentCommand.Depends = strings.Split(content, ",")
			default:
				fmt.Printf("unknown propery %s.\n", parts[0])
			}
		} else if trimmedLine == "" {
			// skipping empty lines
		} else {
			if len(currentCommand.Content) == 0 {
				currentCommand.Content = line
			} else {
				currentCommand.Content = currentCommand.Content + "\n" + line
			}
		}

	}
}

func (c *Config) toLines() []string {
	lines := []string{}
	for _, cm := range c.Commands {
		for _, l := range cm.ToLines() {
			lines = append(lines, l)
		}
	}

	return lines
}

func (c *Config) Print() {
	lines := c.toLines()

	for i, line := range lines {
		fmt.Println(i, line)
	}
}

func (c *Config) Save(path string) error {
	return c.writeLines(c.toLines(), path)
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func (c *Config) readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func (c *Config) writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func (c *Config) AddCommand(cmd *Command) error {
	c.Commands = append(c.Commands, cmd)

	return nil
}

func (c *Config) RmCommand(name string) error {
	cmds := make([]*Command, 0)
	for _, cmd := range c.Commands {
		if cmd.Name != name {
			cmds = append(cmds, cmd)
		}
	}
	c.Commands = cmds

	return nil
}

func (c *Config) HeaderTemplate(name string) []byte {
	x := "# " + name + "\n" +
		"## description: <alias description>\n" +
		"## usage: " + name + " <parameters>\n" +
		"## group: <group>\n" +
		name + "() {\n" +
		"\t# <your code here>\n" +
		"}\n"

	return []byte(x)
}
