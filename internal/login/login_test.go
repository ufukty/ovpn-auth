package login

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ufukty/ovpn-auth/internal/files"
	"gopkg.in/yaml.v3"
)

func loadTestRequestsFile(path string) ([]files.LoginRequest, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer fh.Close()
	c := []files.LoginRequest{}
	err = yaml.NewDecoder(fh).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("decoding requests file: %w", err)
	}
	return c, nil
}

func TestLoginInvalids(t *testing.T) {
	tcs := []string{
		"invalid-totp.yml",
		"invalid-username.yml",
		"invalid-password-totp.yml",
	}

	db, err := files.LoadDatabase("testdata/database.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, loading database: %w", err))
	}

	for _, tn := range tcs {
		t.Run(tn, func(t *testing.T) {
			rqs, err := loadTestRequestsFile(filepath.Join("testdata", tn))
			if err != nil {
				t.Fatal(fmt.Errorf("loading test requests file %q: %w", tn, err))
			}
			for _, rq := range rqs {
				if err := login(db, &rq); err == nil {
					t.Errorf("unexpected success for %q\n", rq.Username)
				} else {
					fmt.Println(err)
				}
			}
		})
	}
}
