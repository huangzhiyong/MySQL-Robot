package common

import (
	"bufio"
	"os"
	"strings"
)

func ReadLinesOffset(filename string, offset int, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var res []string

	r := bufio.NewReader(f)
	for i := 0; i < n+offset || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < offset {
			continue
		}

		res = append(res, strings.Trim(line, "\n"))
	}
	return res, nil
}

func ReadLinesAll(filename string) ([]string, error) {
	return ReadLinesOffset(filename, 0, -1)
}
