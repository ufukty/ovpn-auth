package register

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/alexedwards/argon2id"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/ufukty/ovpn-auth/internal/files"
	"golang.org/x/term"
)

var unavailableUsernames = []files.Username{
	"register",
}

func askUsername() (files.Username, error) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter username: ")
		scanner.Scan()
		username := scanner.Text()
		if err := scanner.Err(); err != nil {
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

	secret := key.String()

	fmt.Println("Copy the secret or scan QR:", secret)
	qrterminal.Generate(secret, qrterminal.H, os.Stdout)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		scanner.Scan()
		nonce := scanner.Text()
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("reading username: %w", err)
		}
		match := totp.Validate(nonce, secret)
		if match {
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
		Memory:      2 ^ 15, // 32 megabytes
		Iterations:  4,
		Parallelism: 2,
		SaltLength:  32,
		KeyLength:   32,
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
	return nil
}
