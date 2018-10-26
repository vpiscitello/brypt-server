# Brypt Server

To build/run:

1) Clone brypt-server into your ~/go/src directory
2) If testing local add the following lines to your hosts file
- 127.0.0.1   access.localhost
- 127.0.0.1   bridge.localhost
- 127.0.0.1   dashboard.localhost
3) Run the following commands in brypt-server/
- $make deps
- $make add_deps
- $make build
4) To run use: $./bin/bserv
