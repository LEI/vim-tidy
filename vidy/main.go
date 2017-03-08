package vidy

import (
	"fmt"
	// "os"
	"strings"
)

var HighLightCommand = "hi"
// HighlightCommands to be ignored
var HighlightCommands = []string{"clear", "link"}
// HighlightGroups keys order
var HighlightGroups = []string{"ctermfg", "ctermbg", "guifg", "guibg", "cterm", "gui"}
// Default number of spaces
var MinIndentSize = 1

// :hi[ghlight] {group-name} {key}={arg}

// Highlight reads a file and reformats `highlight` commands
// Does not handle vim line-continuation
func Highlight(path string) error {
	var lineNr int
	var lines = make([]string, 0)
	var maxLengths = make([]int, len(HighlightGroups)+2)
	var defLines = make([]int, 0) // Line number of highlight definitions
	// Buffer lines
	err := scanFileText(path, func(line string) error {
		lineNr++
		if len(strings.TrimLeft(line, " ")) == 0 {
			// Empty line
			lines = append(lines, line)
			return nil
		}
		fields := strings.Fields(line)
		if len(fields) < 3 || !isHighlightDefinition(fields) {
			// fmt.Fprintf(os.Stderr, "Ignoring line %d: not an highlight group definition\n", lineNr)
			// Not `hi[light] {group-name} ...`
			lines = append(lines, line)
			return nil
		}
		// if len(fields) > len(HighlightGroups) {}
		args := make(map[string]string, 0)
		for _, field := range fields[2:] {
			parts := strings.Split(field, "=")
			if len(parts) != 2 {
				return fmt.Errorf("Invalid line %d in file %s: expecting key/value pair, got '%s'\n", lineNr, path, field)
			}
			args[parts[0]] = parts[1]
		}
		defLines = append(defLines, lineNr)
		hiGroup := highlightGroupMap(fields[1], args)
		for i, field := range strings.Fields(hiGroup) {
			if maxLengths[i] < len(field) {
				maxLengths[i] = len(field)
			}
		}
		lines = append(lines, hiGroup)
		return nil
	})
	if err != nil {
		return err
	}
	// Reformat highligh definitions
	for i, line := range lines {
		var def = false
		for _, n := range defLines {
			if n == i + 1 {
				def = true
				break
			}
		}
		if !def {
			fmt.Println(line)
			continue
		}
		fields := strings.Fields(line)
		// fmt.Println(line, "->", fields)
		if len(fields) == 0 {
			fmt.Println() // FIXME new line
			continue
		}
		if len(fields) != len(HighlightGroups)+2 {
			return fmt.Errorf("invalid highlight group '%s' expected %d fields, got %d", line, len(fields), len(HighlightGroups)+2)
		}
		// var indent = make([]interface{}, len(HighlightGroups)+2)
		// for j, _ := range HighlightGroups {
		// 	indent = append(indent, int(maxLengths[j]) - len(fields[j+2]))
		// }
		// TODO insert old indent
		var str string
		for j, f := range fields {
			str += f
			// Stop at last field
			if j + 1 == len(fields) {
				break
			}
			pad := 1
			if j > 0 {
				pad = maxLengths[j] - len(f)
				// fmt.Fprintf(os.Stderr, "pad: %d - %d + 1 = %d\n", maxLengths[j], len(f), pad)
				if pad < MinIndentSize {
					pad = MinIndentSize
				} else {
					pad += MinIndentSize
				}
			}
			str += strings.Repeat(" ", pad)
		}
		fmt.Println(str)
	}
	// :Tabularize / \+\zs/l0l1
	return nil
}

// Missing keys will be set to `NONE`
func highlightGroupMap(name string, args map[string]string) string {
	str := fmt.Sprintf(HighLightCommand + " %s", name)
	for _, group := range HighlightGroups {
		value := "NONE"
		for k, v := range args {
			if k == group {
				value = v
				break
			}
		}
		str += fmt.Sprintf(" %s=%s", group, value) // ` key=val`
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
