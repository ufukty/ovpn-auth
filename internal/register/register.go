package register

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/ufukty/ovpn-auth/internal/files"
	"golang.org/x/term"
)

var unavailableUsernames = []files.Username{
	"register",
	"validate",
	"version",
}

func ask(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func askUsername() (files.Username, error) {
	for {
		username, err := ask("Enter username: ")
		if err != nil {
			return "", fmt.Errorf("reading username: %w", err)
		}
		if username == "" {
			continue
		}
		return files.Username(username), nil
	}
}

func askPassword() (string, error) {
	for {
		fmt.Print("Enter password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("error reading password: %w", err)
		}
		fmt.Println()
		if len(password) == 0 {
			continue
		}
		return string(password), nil
	}
}

func setTotpSecret(username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ovpn-auth",
		AccountName: username,
	})
	if err != nil {
		return "", fmt.Errorf("generating totp key: %w", err)
	}

	url := key.String()
	fmt.Println("Copy the secret or scan QR:", url)
	qrterminal.Generate(url, qrterminal.M, os.Stdout)

	secret := key.Secret()
	for {
		nonce, err := ask("Enter TOTP nonce: ")
		if err != nil {
			return "", fmt.Errorf("reading totp nonce to validate: %w", err)
		}
		if totp.Validate(nonce, secret) {
			return secret, nil
		}
		fmt.Println("try again")
	}
}

func WithInteraction() error {
	if err := files.CheckDatabase(); err != nil {
		return fmt.Errorf("checking database: %w", err)
	}
	username, err := askUsername()
	if err != nil {
		return fmt.Errorf("asking usename: %w", err)
	}
	if slices.Index(unavailableUsernames, username) != -1 {
		return fmt.Errorf("username is reserved keyword")
	}
	pwd, err := askPassword()
	if err != nil {
		return fmt.Errorf("asking password: %w", err)
	}
	hash, err := argon2id.CreateHash(pwd, &argon2id.Params{
		Memory:      64 * 1024, // 64 megabytes
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  32,
		KeyLength:   128,
	})
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	secret, err := setTotpSecret(string(username))
	if err != nil {
		return fmt.Errorf("setting a totp secret: %w", err)
	}
	db, err := files.LoadDefaultDatabase()
	if err != nil {
		return fmt.Errorf("loading database: %w", err)
	}
	db[username] = files.DatabaseRecord{
		Username:   username,
		Hash:       hash,
		TotpSecret: secret,
	}
	err = db.Save()
	if err != nil {
		return fmt.Errorf("saving changes to database: %w", err)
	}
	fmt.Println("success")
	return nil
}
