#!/bin/bash

docker run -it -d -p 3322:3322 -p 9497:9497 --name immudb codenotary/immudb:latest

