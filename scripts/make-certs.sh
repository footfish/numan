#! /bin/bash 
# Approaches to certs for dev. 
# 1) Use self-signed cert 
# Generate a self-signed cert with openssl
#
# This approach requires the cert to be loaded in the client (or switch off verification)
# in client (numan/main.go) you'll need to load generated cert file with below code. 
#
#	creds, err := credentials.NewClientTLSFromFile(certFile, "")
#	if err != nil {
#		log.Fatalf("cert load error: %s", err)
#	} 
#
#Create self-signed cert 
#openssl req -x509 -newkey rsa:4096 -keyout ../examples/key.pem -out ../examples/cert.pem -addext "subjectAltName = DNS:localhost" -days 365 -nodes  -subj "/CN=*.localhost"

# 2) Use your own trusted CA with minica (preferred)
# installation - https://github.com/jsha/minica/
#
# in client (numan/main.go) you DON'T need to load cert, see below code. 
#
#	creds := credentials.NewTLS(&tls.Config{})
#
# Create localhost cert & trusted authority cert
minica -domains localhost.localdomain,localhost
# copy to target location 
mv localhost.localdomain/*.pem ../examples/
rm -fr localhost.localdomain
# You will now need install minica.pem as trusted on your client computer
# Tips on how to set up - https://gist.github.com/mwidmann/115c2a7059dcce300b61f625d887e5dc
# 
# If it fails get this error "transport: authentication handshake failed: x509: certificate signed by unknown authority"
