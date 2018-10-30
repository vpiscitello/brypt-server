# Brypt Server

To build/run:

1) Clone brypt-server into your ~/go/src directory
2) If testing local: Add the following lines to your hosts file
- 127.0.0.1   access.localhost
- 127.0.0.1   bridge.localhost
- 127.0.0.1   dashboard.localhost
3) If testing local: Generate TLS certificate key pair and place in /config/ssl
- Install "generate_cert" via "go install $GOROOT/src/crypto/tls/generate_cert.go"
    - May need "export $GOBIN=$GOROOT/bin"
- Run "generate_cert -host=127.0.0.1,localhost,access.localhost,bridge.localhost,dashboard.localhost,brypt.com -ca -ecdsa-curve=P384"
- Place "cert.pem" and "key.pem" in /config/ssl
4) Run the following commands in brypt-server/
- $make deps
- $make add_deps
- $make build
5) To run use: $./bin/bserv
6) If testing local: Run "curl -L -k -X GET http://access.localhost:3005"
