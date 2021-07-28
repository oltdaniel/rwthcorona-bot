#!/usr/bin/bash

DOMAIN=$1
BASE=$(pwd $(dirname $0))

if [ "$DOMAIN" == "" ]; then
    read -p "domain: " DOMAIN
fi

if [[ ! -f "$BASE/certificate.key" && ! -f "$BASE/certificate.crt" ]]; then
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certificate.key -out certificate.crt -subj "/CN=$DOMAIN"
fi