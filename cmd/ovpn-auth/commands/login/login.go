// Copyright 2021 Ufuktan Yıldırım (ufukty)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package login

import (
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/pquerna/otp/totp"
	"github.com/ufukty/ovpn-auth/internal/files"
)

func login(users files.Database, r *files.LoginRequest) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unknown error: %s", fmt.Sprint(r))
		}
	}()

	if r.Username == "" {
		return fmt.Errorf("empty field 'username'")
	}
	if r.Password == "" {
		return fmt.Errorf("empty field 'password'")
	}
	if r.TotpNonce == "" {
		return fmt.Errorf("empty field 'totp nonce'")
	}

	user, ok := users[r.Username]
	if !ok {
		return fmt.Errorf("user has not found")
	}

	match, _, err := argon2id.CheckHash(r.Password, user.Hash)
	if err != nil {
		return fmt.Errorf("checking password: %w", err)
	}
	if !match {
		return fmt.Errorf("password mismatch")
	}

	ok = totp.Validate(r.TotpNonce, user.TotpSecret)
	if !ok {
		return fmt.Errorf("totp nonce is invalid")
	}

	return nil
}

func WithFile(path string) error {
	request, err := files.ParseLoginRequest(path)
	if err != nil {
		return fmt.Errorf("reading user input: %w", err)
	}
	db, err := files.LoadDatabase("/etc/openvpn/ovpn_auth_database.yml")
	if err != nil {
		return fmt.Errorf("loading database: %w", err)
	}
	if err := login(db, request); err != nil {
		return err
	}
	return nil
}
