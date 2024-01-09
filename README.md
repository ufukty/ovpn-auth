# ovpn-auth

ovpn-auth is a multi-factor authentication solution for OpenVPN that supports both password and time based one-time-pasword (otp, e.g. Google Authenticator) nonces. It stores passwords after processing them with the state-of-the-art key derivation function Argon2. It is written in Go and has a setup assistance shell script to start using quickly as possible.

## Objectives of the project

If you find the project doesn't meet any of those requirements, create an issue or PR.

-   Support for both password and totp (time-based one time password)
-   Prevent injection
-   Similar time of completion for both valid & invalid requests (one of the possible measures against brute-force & timing attacks)

## Read before use

-   Solutions in this repository may not be safe or secure to use. Review it before use. Take your own risk. If you find an issue, create an issue in GitHub.

-   Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the LICENSE file.

-   Password derivation function takes 32 MiB of space in memory for each login request. So, adjust firewall in a way it will deny abusive amount of requests may originate by attackers from different IPs to take one of the measures against Denial-of-Service attacks.

-   Since, OpenVPN daemon starts the `ovpn-auth` script as the user `nobody`, `secrets.yml` file should be accessible by `nobody`. That means username, password, and otp secret will be able to seen by anyone in the server. While it is not a big problem for argon2 hashes, you should mind the exposure of otp secret.

## Dependencies

Run `make install-dependencies-mac` or `make install-dependencies-linux` to get dependencies.

-   Dependencies of ovpn-auth:
    -   OpenVPN (not included in Makefile)
    -   Argon2 (libargon2)
-   Dependencies of setup assistance:
    -   oathtool (brew/apt)
    -   qr (pip)

## How to use

### Prepare secrets file via setup assistance

Setup assistance will ask you for your choice of username, password and will create a secrets.yml file which contains username, hashed password and otp secret. All you have to do is starting script with command below and following instructions it will show.

```sh
$ wget https://github.com/ufukty/ovpn-auth/raw/main/supplements/setup_assistance.sh
$ bash "setup_assistance.sh"

# if you are on the server
$ mv secrets.yml /etc/openvpn/secrets.yml
# if you are on the development machine
$ scp secrets.yml server:/etc/openvpn/secrets.yml
```

### Download ovpn-auth binary

```sh
# linux-amd64
$ wget https://github.com/ufukty/ovpn-auth/releases/download/v0.2-test/ovpn-auth-210813-linux-amd64.tar.gz
$ tar -xzf ovpn-auth-210813-linux-amd64.tar.gz
$ mv ovpn-auth-210813-linux-amd64.tar.gz /etc/openvpn/ovpn-auth

# macos-amd64
$ wget https://github.com/ufukty/ovpn-auth/releases/download/v0.2-test/ovpn-auth-210813-darwin-amd64.tar.gz
$ tar -xzf ovpn-auth-210813-darwin-amd64.tar.gz
$ mv ovpn-auth-210813-darwin-amd64.tar.gz /etc/openvpn/ovpn-auth
```

### Configure server

```sh
# Locate the `/etc/openvpn` folder of server in your terminal.
$ cd /etc/openvpn

# Adjust permissions and ownership
$ chmod 744 secrets.yml
$ chown root:root secrets.yml
# Caution: That means every user on the server will be able to read the content of secrets file.

$ chmod 755 ovpn-auth
$ chown root:root ovpn-auth

# Enable auth script support for OpenVPN server by editing server.conf file in the server.
$ echo "script-security 2" >>server.conf
$ echo "auth-user-pass-verify /etc/openvpn/ovpn-auth via-file" >>server.conf
```

### Configure clients

In the client configuration, make this update to enable username/password prompt when you try to connect to server:

```sh
$ echo "auth-user-pass" >>/path/to/client_config.ovpn
```

## Login

Enter your password and TOTP code without any character between them when OpenVPN asks. Last 6 digit should be the TOTP code

```sh
$ sudo openvpn client_config.ovpn
...
...
Enter Auth Username:<username>
Enter Auth Password:<password><totp>
...
...
```

## Security Measures

-   Argon2 is used for password hashing with those settings:

    | Setting     | Used value                        | [OWASP suggestion](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#peppering) |
    | ----------- | --------------------------------- | -------------------------------------------------------------------------------------------------------------- |
    | Mode        | id (memory/computation intensive) | id                                                                                                             |
    | Iteration   | 4                                 | 2                                                                                                              |
    | Memory      | 32 MiB                            | 15MiB                                                                                                          |
    | Parallelism | 2                                 | 1                                                                                                              |

## Building from source

### Linux

Tested with Ubuntu 20.04

```sh
cd ~

# Download & install libargon2
wget https://github.com/P-H-C/phc-winner-argon2/archive/refs/tags/20190702.tar.gz
tar -xvf 20190702.tar.gz
cd phc-winner-argon2-20190702
sudo make install
cd -

# Download & install Go
wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
rm go1.16.5.linux-amd64.tar.gz

# Add .../go/bin to PATH environment variable
echo "export PATH=\"\$PATH:/usr/local/go/bin\"" >>~/.bash_profile
source ~/.bash_profile

# Create GO path in your home folder
mkdir -p go/{bin,pkg,src}
GOPATH=~/go

# Download repository
mkdir -p go/src/github.com/ufukty
cd go/src/github.com/ufukty
git clone https://github.com/ufukty/ovpn-auth.git

# Build
cd ovpn-auth
go build .

# Compiled binary has the name 'main'
file main
```

## Test results against timing attacks

Tests are performed in various scenerios and compares average of 50 iterations for both invalid and valid inputs that might come from user/attacker.

| Scenario         | Difference of averages (%) | Difference of averages (ms) | Difference of standart deviations (%) |
| ---------------- | -------------------------- | --------------------------- | ------------------------------------- |
| Invalid username | 3.20 %                     | 0.92 ms                     | 11.28 %                               |
| Empty password   | 99.99 %                    | 23.34 ms                    | 99.99 %                               |
| Invalid password | 0.95 %                     | 0.27 ms                     | 7.59 %                                |

\* In all tests, desired result is `0ms` or `0%`. Results may vary between runs, depends on CPU load coming from other applications runinng on computer.

You can run those tests in your computer with this command:

```sh
go test -v main/validate/...
```

## Contribution

Issues and PRs are welcome.

## Resources

-   https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
-   https://github.com/P-H-C/phc-winner-argon2
-   https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/
-   https://openvpn.net/diy-mfa-setup-community-edition/
-   https://github.com/google/google-authenticator/wiki/Key-Uri-Format

## License

Apache 2.0 License. Read the LICENSE file for information.
