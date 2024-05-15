package files

import (
	"fmt"
	"testing"
)

func Test_InputParse(t *testing.T) {
	tcs := map[string]*LoginRequest{
		"testdata/empty.txt":         nil,
		"testdata/newline.txt":       nil,
		"testdata/no-totp.txt":       nil,
		"testdata/only-pass.txt":     nil,
		"testdata/only-username.txt": nil,
		"testdata/username-totp.txt": nil,
		"testdata/valid.txt": {
			Username:  "username",
			Password:  "passwd",
			TotpNonce: "123456",
		},
		"testdata/zhenyi2.txt": { // https://github.com/ufukty/ovpn-auth/issues/3
			Username:  "zhenyi2",
			Password:  "7758a",
			TotpNonce: "303234",
		},
	}

	for path, expected := range tcs {
		t.Run(path, func(t *testing.T) {
			r, err := ParseLoginRequest(path)
			if (expected == nil) && err == nil {
				t.Fatal("expected to fail, but returned with success")
			} else if (expected != nil) && (err != nil) {
				t.Fatal(fmt.Printf("expected to succeed, but returned with an error: %s", err.Error()))
			} else if expected != nil {
				if expected.Username != r.Username {
					t.Fatal(fmt.Errorf("assert Username, expected %q, got %q", expected.Username, r.Username))
				}
				if expected.Password != r.Password {
					t.Fatal(fmt.Errorf("assert Password, expected %q, got %q", expected.Password, r.Password))
				}
				if expected.TotpNonce != r.TotpNonce {
					t.Fatal(fmt.Errorf("assert TotpNonce, expected %q, got %q", expected.TotpNonce, r.TotpNonce))
				}
			}
		})
	}
}
