#!/bin/bash
# Entry point script for AppImage

SELF=$(readlink -f "$0")
HERE=${SELF%/*}
EXEC="${HERE}/opt/equeselfgo/equeselfgo"

export LD_LIBRARY_PATH="/usr/lib:${HERE}/usr/lib"

exec "${EXEC}"
