#!/bin/bash

# You can edit between these! ----
install_location="/usr/local/opt/libindy"

# -----
# bash exception handling
set -e

get_abs_filename() {
	# $1 : relative filename
	echo "$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
}

install_location=$(get_abs_filename "$install_location")

./tools/prerequisites.sh
./tools/download.sh "$install_location"
./tools/update-deps.sh "$install_location/lib"

# indysdk package tries to find headers from indy/ dir not include/
pushd $install_location
ln -s include indy
popd

# build env loader
cat >env.sh <<EOL
#!/bin/bash

export CGO_CFLAGS="-I""$install_location"
export CGO_LDFLAGS="-L""$install_location""/lib"
EOL

cat >/dev/stdout <<EOF
We have now installed libindy to your given location and generated the
environment variables loading script (env.sh) in this directory. You need to
source the file to set the env variables for CGO.

Add this to your environment files to make it permanent:

	source env.sh

Don't forget to call it!!
EOF

