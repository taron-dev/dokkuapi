package common

import (
	"bufio"
	"os"
)

func readUrls(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, "http://"+scanner.Text())
	}
	return lines, scanner.Err()
}
