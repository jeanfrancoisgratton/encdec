#!/usr/bin/env bash

echo "Installing dependencies";echo
sudo apt-get update && sudo apt update -y
echo;echo;echo "Done. Now installing the Go binaries"
sudo /opt/bin/install_golang.sh 1.23.4 amd64
