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

// Highlight reads a file and reformat `highlight` commands
// Does not handle line continuation
func Highlight(path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	var lineNr int = 0
	scanner := bufio.NewScanner(fi)
LINES:
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()
		// l = strings.TrimLeft(line, " ")
		if len(line) == 0 {
			skip(line)
			continue LINES
		}
		var hiName string
		var hiArgs = make(map[string]string, 0)
		fields := strings.Fields(line)
		for i, word := range fields {
			switch {
			case i == 0 && !strings.HasPrefix(word, "hi"):
				skip(line)
				continue LINES
			case i == 1:
				for _, c := range HighlightCommands {
					if c == word {
						// commandName = word
						skip(line)
						continue LINES
					}
				}
				// Highlight!
			}
			if i > 1 {
				break
			}
		}
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "Ignoring line %d: less than 3 words\n", lineNr)
			skip(line)
			continue LINES
		}
		hiName = fields[1]
		for _, field := range fields[2:] {
			f := strings.Split(field, "=")
			if len(f) != 2 {
				fmt.Fprintf(os.Stderr, "Ignoring line %d: expecting key/value pair, got '%s'\n", lineNr, field)
				skip(line)
				continue
			}
			hiArgs[f[0]] = f[1]
		}
		fmt.Printf("highlight %s", hiName)
		for _, group := range HighlightGroups {
			value := "NONE"
			for k, v := range hiArgs {
				if k == group {
					value = v
					break
				}
			}
			// hi = append(hi, group+"="+value)
			fmt.Printf("%s%s=%s", Separator, group, value)
		}
		fmt.Println() // End of line
	}
	// :Tabularize / \+\zs/l0l1
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func skip(line string) {
	fmt.Println(line)
}
