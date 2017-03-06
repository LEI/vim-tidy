package tidy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HighlightCommands to be ignored
var HighlightCommands = []string{"clear", "link"}
// HighlightGroups defines the keys order
var HighlightGroups = []string{"guifg", "guibg", "gui", "ctermfg", "ctermbg", "cterm", "term"}
// Separator between each key/value pair
var Separator = " "

// :hi[ghlight] {group-name} {key}={arg}

// Highlight reads a file and reformats `highlight` commands
// Does not handle vim line-continuation
func Highlight(path string) error {
	var lineNr int = 0


	scanFile(path, func(line string) error {
		lineNr++
		if len(strings.TrimLeft(line, " ")) == 0 {
			skip(line)
			return nil
		}
		var hiName string
		var hiArgs = make(map[string]string, 0)
		fields := strings.Fields(line)
		if len(fields) < 3 || !isHighlightDefinition(fields) {
			// fmt.Fprintf(os.Stderr, "Ignoring line %d: not an highlight group definition\n", lineNr)
			skip(line)
			return nil
		}
		hiName = fields[1]
		for _, field := range fields[2:] {
			f := strings.Split(field, "=")
			if len(f) != 2 {
				return fmt.Errorf("Invalid line %d in file %s: expecting key/value pair, got '%s'\n", lineNr, path, field)
			}
			hiArgs[f[0]] = f[1]
		}
		hi := HighlightGroup(hiName, hiArgs)
		fmt.Println(hi) // End of line
		return nil
	})
	// :Tabularize / \+\zs/l0l1
	return nil
}

func scanFile(path string, f func(string) error) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		err := f(scanner.Text())
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// HighlightGroup takes a name and a list of key value pairs
// and returns a sorted :highlight command
// Missing keys will be set to `NONE`
func HighlightGroup(name string, args map[string]string) string {
	str := fmt.Sprintf("highlight %s", name)
	for _, group := range HighlightGroups {
		value := "NONE"
		for k, v := range args {
			if k == group {
				value = v
				break
			}
		}
		str += fmt.Sprintf("%s%s=%s", Separator, group, value)
	}
	return str
}

func isHighlightDefinition(fields []string) bool {
	for i, word := range fields {
		switch {
		case i == 0 && !strings.HasPrefix(word, "hi"):
			return false
		case i == 1: // i > 1
			for _, c := range HighlightCommands {
				if c == word {
					return false
				}
			}
			break
		}
	}
	return true
}

func skip(line string) {
	fmt.Println(line)
}
