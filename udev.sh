#!/bin/bash
# This script sets up udev rules for Neewer devices to ensure proper permissions and access.
# Must be run as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

# Get the actual user (not root)
if [ -n "$SUDO_USER" ]; then
    USER="$SUDO_USER"
else
    USER=$(logname 2>/dev/null || echo "$USER")
fi

# Detect serial group by checking what exists
if getent group dialout > /dev/null 2>&1; then
    DEFAULT_GROUP="dialout"
elif getent group uucp > /dev/null 2>&1; then
    DEFAULT_GROUP="uucp"
elif getent group plugdev > /dev/null 2>&1; then
    DEFAULT_GROUP="plugdev"
else
    echo "Warning: No standard serial group found. Creating 'dialout' group."
    DEFAULT_GROUP="dialout"
fi

DEFAULT_VENDOR_ID="1a86"
DEFAULT_PRODUCT_ID="7523"
DEFAULT_UDEV_RULES_FILE="/etc/udev/rules.d/99-neewerctl.rules"
DEFAULT_MODE="0666"

VENDOR_ID=$1
PRODUCT_ID=$2
UDEV_RULES_FILE=$3
MODE=$4
GROUP=$5

# Set defaults if parameters are not provided
VENDOR_ID=${VENDOR_ID:-$DEFAULT_VENDOR_ID}
PRODUCT_ID=${PRODUCT_ID:-$DEFAULT_PRODUCT_ID}
UDEV_RULES_FILE=${UDEV_RULES_FILE:-$DEFAULT_UDEV_RULES_FILE}
MODE=${MODE:-$DEFAULT_MODE}
GROUP=${GROUP:-$DEFAULT_GROUP}

echo "Setting up udev rules for Neewer devices:"
echo "  User: $USER"
echo "  Vendor ID: $VENDOR_ID"
echo "  Product ID: $PRODUCT_ID"
echo "  Udev Rules File: $UDEV_RULES_FILE"
echo "  Mode: $MODE"
echo "  Group: $GROUP"
echo ""
read -p "Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Udev setup cancelled."
    exit 1
fi

# Create group if it doesn't exist (only for non-standard groups)
if ! getent group "$GROUP" > /dev/null 2>&1; then
    echo "Creating group $GROUP..."
    groupadd "$GROUP"
fi

# Add user to group if not already a member
if ! groups "$USER" 2>/dev/null | grep -q "\b$GROUP\b"; then
    echo "Adding user $USER to group $GROUP..."
    usermod -aG "$GROUP" "$USER"
    NEEDS_LOGOUT=true
else
    echo "User $USER is already a member of group $GROUP."
    NEEDS_LOGOUT=false
fi

# Create udev rules file
echo "Creating udev rules file at $UDEV_RULES_FILE..."
cat <<EOL > "$UDEV_RULES_FILE"
# Neewer devices - allow user access via $GROUP group
SUBSYSTEM=="tty", ATTRS{idVendor}=="$VENDOR_ID", ATTRS{idProduct}=="$PRODUCT_ID", MODE="$MODE", GROUP="$GROUP"
EOL

# Reload udev rules
echo "Reloading udev rules..."
udevadm control --reload-rules
udevadm trigger

echo ""
echo "✓ Udev rules installed successfully!"
echo ""
echo "IMPORTANT: Please complete these steps:"
if [ "$NEEDS_LOGOUT" = true ]; then
    echo "  1. Log out and log back in (group membership takes effect)"
fi
echo "  2. Unplug and replug your Neewer device"
echo "  3. Test with: neewerctl status"
echo ""
