#!/bin/bash

# indysdk package tries to find headers from indy/ dir not include/
ln -s ./include indy 

export CGO_CFLAGS="-I""$PWD"
export CGO_LDFLAGS="-L""$PWD""/lib"

