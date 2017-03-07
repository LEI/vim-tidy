package vidy

import (
	"bufio"
	"os"
)

func scanFileText(path string, f func(string) error) error {
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
