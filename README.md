# `ovpn-auth`

`ovpn-auth` is a easy-to-use multi-factor authentication solution for OpenVPN that supports both password and Time-Based OTP nonces.

## Features

- 2FA login with username, password and TOTP nonce
- Multiple users
- Built-in mechanism to register. Scan QR code printed on terminal to TOTP app e.g. Google Authenticator
- No dependencies & single binary => easy to deploy & use
- Passwords are stored as Argon2id hashes

### Register new users

```sh
$ ovpn-auth register

Enter username: ufukty
Enter password:
Copy the secret or scan QR: otpauth://totp/ovpn-auth:ufukty?algorithm=SHA1&digits=6&issuer=ovpn-auth&period=30&secret=V72ILSVNVZ2FMGMJKI5UPC64XXFGPJXK
# QR code will be printed here
Enter TOTP nonce: 241988
success
```

This will ask username, password and print a totp secret with QR code for it. It needs `/etc/openvpn` directory to exist.

### Login

> [!IMPORTANT]
> Append 6 digit TOTP nonce just to the end of your password.

```sh
$ sudo openvpn client_config.ovpn
...
Enter Auth Username:<username>
Enter Auth Password:<password><totp>
```

## How to use

### Configure server

```sh
# Locate the `/etc/openvpn` folder of server in your terminal.
cd /etc/openvpn

# Adjust permissions and ownership
chmod 744 ovpn_auth_database.yml # Read caution below
chown root:root ovpn_auth_database.yml

chmod 755 ovpn-auth
chown root:root ovpn-auth

# Enable auth script support for OpenVPN server by editing server.conf file in the server.
echo "script-security 2" >>server.conf
echo "auth-user-pass-verify /etc/openvpn/ovpn-auth via-file" >>server.conf
```

> [!CAUTION]
> Since, OpenVPN daemon starts the `ovpn-auth` as the user `nobody`, also the `ovpn_auth_database.yml` file should be accessible by `nobody`. That means username, password, and TOTP secret will be able to seen by anyone in the server. While it is not a big problem for argon2 hashes, you should mind the exposure of otp secret.

### Configure clients

In the client configuration, make this update to enable username/password prompt when you try to connect to server:

```sh
$ echo "auth-user-pass" >>/path/to/client_config.ovpn
```

### Deployment

Database file is portable. Just be sure it has correct permissions to let user `nobody` can read.

```yml
# /etc/openvpn/ovpn_auth_database.yml
- username: ufukty
  hash: $argon2id$v=19$m=13,t=2,p=1$fdxUDnZZ9kPXRUQF8Z3DP1OkwVOw1bWwEB4y/C3RrQ8$x2D/sH/7MlWMrz4Oc2FgMR5pqnw1ZbVSJY2oRpTc1bthOJtRLRAt9IRpV/XZaxSg/8q6ewnz2X2igtOy48uHkuFkRNmIIjpzAetmSa5cqKJWxT1iAA+XNfS54+WqvdjWY4uvi8jKwqaBJupsGWPjwMi/JGUip4mu2LtkQPDPQ5I
  totp-secret: GCH35RULQC42SF2BGXCS2PZFY3KWN3XV
```

## Security

### Argon2 configuration

| Parameter       | Value  |
| --------------- | ------ |
| Mode            | id     |
| Memory (m)      | 64 MiB |
| Iterations (t)  | 2      |
| Parallelism (p) | 1      |
| Salt            | 32     |
| Key             | 128    |

> [!CAUTION]
> Beware that since hashing needs significant amount of RAM to be utilized for each request, external adjustments to block requests that might lead to a DoS attack may be necessary.

### Timing attacks

```sh
$ /usr/local/go/bin/go test -timeout 10s -run ^TestTimings$ github.com/ufukty/ovpn-auth/internal/login -v -count=1

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

Apache 2.0 License.

> [!CAUTION]
> Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> Read the LICENSE file for information.
