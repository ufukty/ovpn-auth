#!/bin/bash

# Copyright 2021 Ufuktan Yıldırım (ufukty)
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
# http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

colr_lre="\033[91m" # l_red
colr_lye="\033[93m" # l_yellow
colr_lma="\033[95m" # l_magenta
colr_lgr="\033[92m" # l_green
colr_rst_all="\033[0m"
colr_rst_bold_bright="\033[1m\033[21m"

function log_err() {
    echo -e "${colr_rst_bold_bright}${colr_lre}[ERROR] $1${colr_rst_all}" >&2
}

function log_ntc() {
    echo -e "${colr_rst_bold_bright}${colr_lye}[NOTICE] $1${colr_rst_all}" >&2
}

function log_inf() {
    echo -e "${colr_rst_bold_bright}${colr_lma}[INFO] $1${colr_rst_all}"
}

function log_suc() {
    echo -e "${colr_rst_bold_bright}${colr_lgr}[SUCCESS] $1${colr_rst_all}" >&2
}

if [[ -e "secrets.yml" ]]; then
    log_err "There is already a secrets.yml file in this folder. Aborting."
    exit 1
fi

username=""
password_clear=""
password_salt="$(head -n 1024 /dev/urandom | sha512sum | cut -b 1-32)"
password_hash_computed=""
otp_secret="$(head -n 100 /dev/urandom | base32 | cut -b 1-32 | head -n 1)"

log_inf "Checking if argon2 is available..."
if [[ ! $(which argon2) ]]; then
    log_err "argon2 is not available in \$PATH, follow instructions https://github.com/P-H-C/phc-winner-argon2"
    exit 2
fi

read -p "> Type your username: " username

while [[ "1" ]]; do
    read -p "> Type your password: " -s password_clear && echo ""
    read -p "> Re-type your password: " -s password_clear_retype && echo ""
    if [[ "$password_clear" == "$password_clear_retype" ]]; then
        log_inf "Passwords are matching."
        break
    else
        log_ntc "Passwords are not matching. Try again."
    fi
done

log_inf "Running argon2, it can take couple minutes..."

password_hash_computed="$(echo -n "$password_clear" | argon2 "$password_salt" -id -t 100 -m 17 -p 2 -e)"

# FIXME: is it necessary?
unset password_clear

cat >secrets.yml <<HERE
username: $username
salt: $password_salt
hash: $password_hash_computed
otp_secret: $otp_secret
HERE

log_inf "secret.yml is created"

if [[ "$EUID" > 0 && -d /etc/openvpn/ ]]; then
    log_inf "Script is running by root user and '/etc/openvpn' does exists.
    So, script will adjust permission and ownership of file and move it to /etc/openvpn behalf of you..."
    chmod 400 "secrets.yml"
    chown root:root "secrets.yml"
    mv "secrets.yml" "/etc/openvpn/secrets.yml"
else
    log_ntc "Script isn't running as root or /etc/openvpn doesn't exists."
    log_ntc "So, you have to adjust A) permission to 400, B) ownership to root, C) move the file into /etc/openvpn/secrets.yml"
fi

echo "Region information will be embedded into the Authenticator string."
read -p "> Type the region (eg. aws-fra): " issuer_name

log_inf "Open this link with internet browser in your phone, then it should redirect you to Authenticator app:"
echo "otpauth://totp/OpenVPN:${username}@${issuer_name}?secret=${otp_secret}&issuer=OpenVPN"

while [[ "1" ]]; do
    read -p "> Type the OTP code produced by the authenticator: " input_otp_nonce
    computed_nonce="$(oathtool --totp --base32 "$otp_secret")"
    if [[ "$input_otp_nonce" =~ ^${computed_nonce}$ ]]; then break; fi
done

log_suc "Script should be finished succesfully."
