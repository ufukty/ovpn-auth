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
	"fmt"
	"os"

	"github.com/ufukty/ovpn-auth/internal/login"
	"github.com/ufukty/ovpn-auth/internal/register"
)

var Version = ""

func dispatch() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("not enough argument")
	}
	switch arg := os.Args[1]; arg {
	case "register":
		if err := register.WithInteraction(); err != nil {
			return fmt.Errorf("performing registration interaction: %w", err)
		}
	case "validate":
		if err := login.WithInteraction(); err != nil {
			return fmt.Errorf("performing login validation flow: %w", err)
		}
	case "version":
		fmt.Println(Version)
	default:
		if err := login.WithFile(arg); err != nil {
			return fmt.Errorf("checking login: %w", err)
		}
	}

	return nil
}

func main() {
	if err := dispatch(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
