#!/bin/bash

# You can edit between these! ----
install_location="./libindy2"

# -----
# bash exception handling
set -e

. ./tools/prerequisites.sh
. ./tools/download.sh "$install_location"
. ./tools/update-deps.sh "$install_location/lib"
