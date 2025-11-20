package files

import (
	"bufio"
	"fmt"
	"os"
)

type LoginRequest struct {
	Username  Username `yaml:"username"`
	Password  string   `yaml:"password"`
	TotpNonce string   `yaml:"totp-nonce"`
}

func ParseLoginRequest(filename string) (*LoginRequest, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("reading user input file: %w", err)
	}
	defer fh.Close() //nolint:errcheck

	var lines []string
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("the file passed by OpenVPN doesn't contain 2 lines of content")
	}
	if len(lines[0]) < 1 {
		return nil, fmt.Errorf("username is shorther than 1 character")
	}
	if len(lines[1]) < 7 {
		return nil, fmt.Errorf("concatenated password and TOTP nonce is shorther than 7 characters, probably missing one of those")
	}

	return &LoginRequest{
		Username:  Username(lines[0]),
		Password:  lines[1][:len(lines[1])-6],
		TotpNonce: lines[1][len(lines[1])-6:],
	}, nil
}
