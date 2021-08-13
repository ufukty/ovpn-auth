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

package file

import (
	"bufio"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type StoredSecrets struct {
	Username   string
	Hash       string
	Otp_secret string
}

type UserInput struct {
	Username  string
	Password  string
	Otp_nonce string
}

func LoadSecrets(filename string) StoredSecrets {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	secrets := StoredSecrets{}

	err = yaml.Unmarshal(content, &secrets)
	if err != nil {
		log.Fatal(err)
	}

	return secrets
}

func LoadUserInput(filename string) UserInput {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return UserInput{
		Username:  lines[0],
		Password:  lines[1][:len(lines[1])-6],
		Otp_nonce: lines[1][len(lines[1])-6:],
	}
}
