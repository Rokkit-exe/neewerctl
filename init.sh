#!/bin/bash

# This script initializes the development environment by setting up necessary configurations and dependencies.

#must be run as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# making sure we are in the right directory
if [ "$(basename "$PWD")" != "neewerctl" ]; then
  echo "Please run this script from the neewerctl directory."
  exit 1
fi

# Remove service file if it exists
if [ -f /etc/systemd/system/neewerd.service ]; then
    echo "Removing existing neewerd.service file..."
    systemctl disable neewerd.service
    rm /etc/systemd/system/neewerd.service
fi

# Remove executable file if it exists
if [ -f /usr/local/bin/neewerctl ]; then
    echo "Removing existing neewerctl executable..."
    rm /usr/local/bin/neewerctl
fi

# Building binary
if [ ! -d ./bin ]; then
  mkdir ./bin
fi
echo "Building neewerctl binary..."
go build -o ./bin/neewerctl main.go
cp ./bin/neewerctl /usr/local/bin/neewerctl


# Setting up config file
echo "Setting up configuration file to $HOME/.config/neewerctl/config.yaml..."
mkdir -p "$HOME/.config/neewerctl"
cp config.yaml "$HOME/.config/neewerctl/config.yaml"
chmod -R 755 "$HOME/.config/neewerctl"

printf "\n\n"
echo "neewerctl binary location:"
echo "  - /usr/local/bin/neewerctl"
echo "Config location:"
echo "  - $HOME/.config/neewerctl/config.yaml"
echo
echo "Initialization complete."

