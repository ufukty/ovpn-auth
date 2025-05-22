package login

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/pquerna/otp/totp"
	"github.com/ufukty/ovpn-auth/internal/files"
)

func ask(prompt string) (string, error) {
	for {
		fmt.Printf("Enter %s:", prompt)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", prompt, err)
		}
		input = strings.TrimSuffix(input, "\n")
		if input == "" {
			continue
		}
		return input, nil
	}
}

type Args struct {
	Database string
}

func WithInteraction() error {
	args := &Args{}
	flag.StringVar(&args.Database, "database", "/etc/openvpn/ovpn_auth_database.yml", "Relative or absolute path of Ovpn-Auth database that will be used for writing the credentials. If there is no file in given path, Ovpn-Auth will create a new database file.")
	flag.Parse()

	db, err := files.LoadDatabase(args.Database)
	if err != nil {
		return fmt.Errorf("loading database: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unknown error: %s", fmt.Sprint(r))
		}
	}()

	username, err := ask("username")
	if err != nil {
		return err
	}

	user, ok := db[files.Username(username)]
	if !ok {
		return fmt.Errorf("user has not found")
	}

	fmt.Println("Note, password will appear below as you type. In validation mode, TOTP nonce will be asked later.")
	password, err := ask("password")
	if err != nil {
		return err
	}

	match, _, err := argon2id.CheckHash(password, user.Hash)
	if err != nil {
		return fmt.Errorf("checking password: %w", err)
	}
	if !match {
		return fmt.Errorf("password mismatch")
	}

	totpNonce, err := ask("TOTP nonce")
	if err != nil {
		return err
	}

	ok = totp.Validate(totpNonce, user.TotpSecret)
	if !ok {
		return fmt.Errorf("totp nonce is invalid")
	}

	fmt.Println("Success, credentials are correct.")
	return nil
}
