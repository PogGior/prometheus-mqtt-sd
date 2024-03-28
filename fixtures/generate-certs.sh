#!/bin/sh

# Save current directory in variable
script_path=$(readlink -f "$0")
root_dir=$(dirname "$script_path")

# Generate CA private key
openssl genrsa -out $root_dir/ca.key 2048

# Generate self-signed CA certificate
openssl req -x509 -new -nodes -key $root_dir/ca.key -subj "/CN=example.com" -days 3650 -out $root_dir/ca.crt

# Generate client private key
openssl genrsa -out $root_dir/client.key 2048

# Generate client certificate signing request (CSR)
openssl req -new -key $root_dir/client.key -subj "/CN=client.example.com" -out $root_dir/client.csr

# Use the CA certificate and key to sign the client CSR and get a client certificate
openssl x509 -req -in $root_dir/client.csr -CA $root_dir/ca.crt -CAkey $root_dir/ca.key -CAcreateserial -out $root_dir/client.crt -days 3650

# Clean up
rm $root_dir/client.csr
rm $root_dir/ca.key
rm $root_dir/ca.srl