# `ovpn-auth`

`ovpn-auth` is a easy-to-use multi-factor authentication solution for OpenVPN that supports both password and Time-Based OTP nonces.

## Features

- 2FA login with username, password and TOTP nonce
- Multiple users
- Built-in mechanism to register. Scan QR code printed on terminal to TOTP app e.g. Google Authenticator
- No dependencies & single binary => easy to deploy & use
- Passwords are stored as Argon2id hashes

> [!CAUTION]
> Solutions in this repository may not be safe or secure to use. Review it before use. Take your own risk. If you find an issue, report.

> [!CAUTION]
> Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the LICENSE file.

> [!CAUTION]
> Password derivation function takes 32 MiB of space in memory for each login request. So, adjust firewall in a way it will deny abusive amount of requests may originate by attackers from different IPs to take one of the measures against Denial-of-Service attacks.

> [!CAUTION]
> Since, OpenVPN daemon starts the `ovpn-auth` script as the user `nobody`, `ovpn_auth_database.yml` file should be accessible by `nobody`. That means username, password, and TOTP secret will be able to seen by anyone in the server. While it is not a big problem for argon2 hashes, you should mind the exposure of otp secret.

## How to use

### Configure server

```sh
# Locate the `/etc/openvpn` folder of server in your terminal.
$ cd /etc/openvpn

# Adjust permissions and ownership
$ chmod 744 ovpn_auth_database.yml
$ chown root:root ovpn_auth_database.yml
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

## Register

```sh
ovpn-auth register
```

This will ask username, password and print a totp secret with QR code for it. It needs `/etc/openvpn` directory to exist.

## Login

Append 6 digit TOTP nonce just to the end of your password.

```sh
$ sudo openvpn client_config.ovpn
...
Enter Auth Username:<username>
Enter Auth Password:<password><totp>
```

## Security Measures

- Argon2 is used for password hashing with those settings:

  | Setting     | `ovpn-auth` | [OWASP suggestion](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id) |
  | ----------- | ----------- | ------------------------------------------------------------------------------------------------------------- |
  | Mode        | id\*        | id\*                                                                                                          |
  | Iteration   | 4           | 2                                                                                                             |
  | Memory      | 32 MiB      | 15MiB                                                                                                         |
  | Parallelism | 2           | 1                                                                                                             |

_\* id = both memory and compute intensive_

## Test results against timing attacks

```sh
/usr/local/go/bin/go test -timeout 10s -run ^TestTimings$ github.com/ufukty/ovpn-auth/internal/login -v -count=1

=== RUN   TestTimings
=== RUN   TestTimings/invalid-totp.yml
=== RUN   TestTimings/invalid-username.yml
=== RUN   TestTimings/invalid-password-totp.yml
total durations for all requests in same test set:
    invalid-totp.yml          => 3.296492025s
    invalid-username.yml      => 11.827µs
    invalid-password-totp.yml => 11.643µs
average durations per request:
    invalid-totp.yml          => 66ms
    invalid-username.yml      => 0ms
    invalid-password-totp.yml => 0ms
standard deviations (amongst all requests in one set, individually):
    invalid-totp.yml          => 2.80ms
    invalid-username.yml      => 0.00ms
    invalid-password-totp.yml => 0.00ms
--- PASS: TestTimings (3.30s)
    --- PASS: TestTimings/invalid-totp.yml (3.30s)
    --- PASS: TestTimings/invalid-username.yml (0.00s)
    --- PASS: TestTimings/invalid-password-totp.yml (0.00s)
PASS
ok      github.com/ufukty/ovpn-auth/internal/login      3.678s
```

## Contribution

Issues and PRs are welcome.

## Resources

- https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
- https://github.com/P-H-C/phc-winner-argon2
- https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/
- https://openvpn.net/diy-mfa-setup-community-edition/
- https://github.com/google/google-authenticator/wiki/Key-Uri-Format

## License

Apache 2.0 License. Read the LICENSE file for information.
