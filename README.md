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

# Optionnal
sudo cp neewerctl /usr/local/bin/neewerctl
```

## Usage

```bash
# must be run as root
# default port: /dev/ttyUSB0
sudo neewerctl --port "/dev/ttyUSB0" --brightness [0-100] --temperature [2700-7000]
sudo neewerctl --on
sudo neewerctl --off
sudo neewerctl --profile [cold|sunlight|afternoon|sunset|candle]
```

## Find Device

```bash
lsusb
# Output
# Bus 003 Device 012: ID 1a86:7523 QinHeng Electronics CH340 serial converter
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
