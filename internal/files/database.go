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

package files

import (
	"fmt"
	"os"

	"github.com/ufukty/ovpn-auth/internal/utils"
	"gopkg.in/yaml.v3"
)

const dbpath = "/etc/openvpn/ovpn_auth_database.yml"

type Username string

type DatabaseRecord struct {
	Username   Username `yaml:"username"`
	Hash       string   `yaml:"hash"`
	TotpSecret string   `yaml:"totp-secret"`
}

type Database map[Username]DatabaseRecord

func LoadDatabase(path string) (Database, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening database file: %w", err)
	}
	defer fh.Close()
	dbrcs := []DatabaseRecord{}
	err = yaml.NewDecoder(fh).Decode(&dbrcs)
	if err != nil {
		return nil, fmt.Errorf("decoding database: %w", err)
	}
	db := utils.Mapify(dbrcs, func(dbrc DatabaseRecord) Username { return dbrc.Username })
	return Database(db), nil
}

func LoadDefaultDatabase() (Database, error) {
	return LoadDatabase(dbpath)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CheckDatabase() error {
	if !fileExists(dbpath) {
		fh, err := os.OpenFile(dbpath, os.O_CREATE|os.O_WRONLY, 0700)
		if err != nil {
			return fmt.Errorf("creating database file: %w", err)
		}
		fh.Close()
	}
	return nil
}

func (db Database) Save() error {
	fh, err := os.OpenFile(dbpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		return fmt.Errorf("opening database file: %w", err)
	}
	err = yaml.NewEncoder(fh).Encode(db)
	if err != nil {
		return fmt.Errorf("writing to database file: %w", err)
	}
	return nil
}
