# Producing test data

**Producing passwords**

Aggregation of:

-   Top used 25 password
-   Random passwords created with character range [A-Za-z0-9]
-   Random passwords created with character range [!-~]

**Producing salt and hash values**

```sh
cat test_data_passwords.yml | while read -r PASS; do
    SALT="$(head -n 1024 /dev/urandom | sha512sum | cut -b 1-32)"
    HASH="$(echo -n "$PASS" | argon2 "$SALT" -id -t 4 -m 15 -p 2 -e)"
    echo "- password: $PASS"
    echo "  salt: $SALT"
    echo "  hash: $HASH"
done
```

**Producing TOTP secrets**

```sh
head -n 100 /dev/urandom | base32 | cut -b 1-32 | head -n 50
```
