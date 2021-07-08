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

package main

import (
	"log"

	"github.com/pquerna/otp/totp"
	"github.com/tvdburgt/go-argon2"
)

type StoredSecrets struct {
	Username   string
	Salt       string
	Hash       string
	Otp_secret string
}

type UserInput struct {
	Username  string
	Password  string
	Otp_nonce string
}

func validateUsername(storedUsername string, inputUsername string) bool {
	return storedUsername == inputUsername
}

func validatePassword(storedHash string, storedSalt string, inputPassword string) bool {
	result, err := argon2.VerifyEncoded(storedHash, []byte(inputPassword))
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func validateOTP(storedOtpSecret string, inputOtpNonce string) bool {
	return totp.Validate(inputOtpNonce, storedOtpSecret)
}
