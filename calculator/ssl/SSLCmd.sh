# Server common name
SERVER_CN=localhost

# Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:123456 -des3 -out ca.key 4096
openssl req -passin pass:123456 -new -x509 -days 3650 -key ca.key -out ca.crt
