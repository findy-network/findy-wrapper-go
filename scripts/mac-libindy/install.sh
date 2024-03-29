#!/bin/bash

. tools/functions.sh

install_location=${1:-/usr/local/opt/libindy}

if [[ "$1" != "" ]]; then
	read -p "Please confirm the installation location [${install_location}]: " location
	install_location=${location:-$install_location}
fi

install_location=$(get_abs_filename "$install_location")
echo "Installing to: ""$install_location"
echo "This will take a moment. Please wait..."

# bash "exception handling"
set -e

./tools/prerequisites.sh
./tools/download.sh "$install_location"
./tools/update-deps.sh "$install_location/lib"

# indysdk package tries to find headers from indy/ dir not include/
pushd $install_location > /dev/null
ln -s include indy
popd > /dev/null

# build env loader
cat >"$install_location"/env.sh <<EOF
#!/bin/bash

export CGO_CFLAGS="-I""$install_location"
export CGO_LDFLAGS="-L""$install_location""/lib"
EOF

cat >/dev/stdout <<EOF

Congrulations!

We have now installed libindy to your system and generated the environment
variables loading script (env.sh) into the installation directory. You need to
source the file to set the env variables for CGO by yourself.

Add this to your environment files to make it permanent:

	source ${install_location}/env.sh

Don't forget to call it for this shell session as well. It's in your clipboard.
EOF

printf "source %s/env.sh" "$install_location" | pbcopy

