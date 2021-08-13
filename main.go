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
	"main/file"
	"main/validate"
	"os"
)

func main() {
	userInput := file.LoadUserInput(os.Args[1])
	storedSecrets := file.LoadSecrets("secrets.yml")

	if validate.Final(storedSecrets, userInput) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
