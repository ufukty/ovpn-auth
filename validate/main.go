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

package validate

import (
	"log"
	"main/file"

	"github.com/pquerna/otp/totp"
	"github.com/tvdburgt/go-argon2"
)

func Username(storedUsername string, inputUsername string) bool {
	return storedUsername == inputUsername
}

func Password(storedHash string, inputPassword string) bool {
	if inputPassword == "" {
		return false
	}
	result, err := argon2.VerifyEncoded(storedHash, []byte(inputPassword))
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func OTP(storedOtpSecret string, inputOtpNonce string) bool {
	return totp.Validate(inputOtpNonce, storedOtpSecret)
}

func Final(storedSecrets file.StoredSecrets, userInput file.UserInput) bool {
	// Perform both OTP and Password checks before sending the response.
	// It is a measure against timing-attack.
	valid_usn := Username(storedSecrets.Username, userInput.Username)
	valid_otp := OTP(storedSecrets.Otp_secret, userInput.Otp_nonce)
	valid_pwd := Password(storedSecrets.Hash, userInput.Password)

	if valid_usn && valid_otp && valid_pwd {
		return true
	} else {
		return false
	}
}
