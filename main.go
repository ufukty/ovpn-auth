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
	"os"
)

func main() {
	userInput := loadUserInput(os.Args[1])
	storedSecrets := loadSecrets("secrets.yml")

	// Perform both OTP and Password checks before sending the response.
	// It is a measure against timing-attack.
	valid_usn := validateUsername(storedSecrets.Username, userInput.Username)
	valid_otp := validateOTP(storedSecrets.Otp_secret, userInput.Otp_nonce)
	valid_pwd := validatePassword(storedSecrets.Hash, storedSecrets.Salt, userInput.Password)

	if valid_usn && valid_otp && valid_pwd {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
