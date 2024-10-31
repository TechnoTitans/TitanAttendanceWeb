package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	name = strings.TrimSpace(name)
	return name, nil
}
