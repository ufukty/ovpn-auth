package files

import (
	"fmt"
	"testing"
)

func Test_InputParse(t *testing.T) {
	r, err := ParseLoginRequest("assets/example.txt")
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}

	if r.Username != "ufukty" {
		t.Errorf("assert, username expected to be %q got %q", "ufukty", r.Username)
	}

	if r.Password != "Hello World!" {
		t.Errorf("assert, password expected to be %q got %q", "Hello World!", r.Password)
	}

	if r.TotpNonce != "111111" {
		t.Errorf("assert, totpNonce expected to be %q got %q", "111111", r.TotpNonce)
	}
}
