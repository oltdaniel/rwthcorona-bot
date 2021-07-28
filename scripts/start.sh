#!/usr/bin/bash

set +x

BASE=$(pwd $(dirname $0))

if [ -f "$BASE/.env" ]
then
    source "$BASE/.env"
else
    echo "ensure you have an .env file"
fi

if [[ ! -f $BASE/certificate.key || ! -f $BASE/certificate.csr ]]; then
    $BASE/scripts/certificate.sh "$HOSTNAME"
fi

TELEGRAM_TOKEN="$TELEGRAM_TOKEN" CERTIFICATE_FILE="$CERTIFICATE_FILE" KEY_FILE="$KEY_FILE" HOSTNAME="$HOSTNAME" go run main.go