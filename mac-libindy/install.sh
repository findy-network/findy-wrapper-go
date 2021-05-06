#!/bin/bash

# You can edit between these! ----
install_location="./libindy"

# -----
# bash exception handling
set -e

. ./prerequisites.sh
. ./download.sh "$install_location"
. ./update-deps.sh "$install_location"
