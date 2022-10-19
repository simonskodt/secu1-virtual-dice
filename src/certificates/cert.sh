#!/bin/bash

for name in "alice" "bob"; do
    openssl req \
        -x509 \
        -newkey rsa:4096 \
        -keyout "$name.key.pem" \
        -out "$name.cert.pem" \
        -sha256 \
        -days 365 \
        -nodes \
        -subj "/CN=$name" \
        -addext "subjectAltName = DNS:localhost,DNS:$name"
done