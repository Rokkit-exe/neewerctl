# neewerctl

### Unofficial command-line tool to control Neewer smart lights

## Features

- Turn Neewer smart lights on and off.
- Adjust brightness levels.
- Change color temperature.
- Set predefined lighting profiles.

## Supported Devices

- Neewer PL81

> [!NOTE] Serial
> `Neewerctl` has only been tested on the `Neewer PL81` but might work with other `Neewer lights` that are using `serial`.

## Requirements

- Go

## Installation

```bash
git clone https://github.com/Rokkit-exe/neewerctl.git
cd neewerctl

go build -o neewerctl main.go

# Optionnal: move binary to /usr/local/bin for global usage
sudo cp neewerctl /usr/local/bin/neewerctl
```

## Usage

```bash
# must be run as root
# default port: /dev/ttyUSB0

# Adjust brightness (0-100) and temperature (2700-7000K)
sudo neewerctl set --device "/dev/ttyUSB0" --brightness [0-100] --temperature [2700-7000]

# Turn light on or off
sudo neewerctl on -d "/dev/ttyUSB0"
sudo neewerctl off -d "/dev/ttyUSB0"

# Set predefined profile
sudo neewerctl set --profile [cold|sunlight|afternoon|sunset|candle] -d "/dev/ttyUSB0"

# List connected Neewer devices
sudo neewerctl list

# list specific device information
sudo neewerctl list -d "/dev/ttyUSB0"
```

## Find Device

```bash
lsusb
# Output
# Bus xxx Device xxx: ID 1a86:7523 QinHeng Electronics CH340 serial converter
```

## Read Serial Port

1. Run these two command to read the serial port

```bash
sudo stty -F /dev/ttyUSB0 raw -echo 115200
sudo cat /dev/ttyUSB0 > neewer_dump.bin
```

2. Adjust brightness/temperature with physical button on the neewer light
3. `CTRL+C` to stop reading the serial port
4. Run this command to analyse the output

```bash
hexdump -C neewer_dump.bin
```

## Profiles

These profiles are the original presets from the `Neewer` app:

- `cold`: Brightness `100%`, Temperature `7000K`
- `sunlight`: Brightness `28%`, Temperature `5600K`
- `afternoon`: Brightness `16%`, Temperature `5000K`
- `sunset`: Brightness `16%`, Temperature `4500K`
- `candle`: Brightness `28%`, Temperature `3400K`
