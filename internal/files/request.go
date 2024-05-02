package files

import (
	"bufio"
	"fmt"
	"os"
)

type LoginRequest struct {
	Username  Username
	Password  string
	TotpNonce string
}

func ParseLoginRequest(filename string) (*LoginRequest, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("reading user input file: %w", err)
	}
	defer fh.Close()

	var lines []string
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return &LoginRequest{
		Username:  Username(lines[0]),
		Password:  lines[1][:len(lines[1])-6],
		TotpNonce: lines[1][len(lines[1])-6:],
	}, nil
}
