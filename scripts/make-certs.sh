#! /bin/bash 
# Creates local certs
openssl req -x509 -newkey rsa:4096 -keyout ../examples/server-key.pem -out ../examples/server-cert.pem -addext "subjectAltName = DNS:localhost" -days 365 -nodes 