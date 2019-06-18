#!/bin/sh

NAME=$1

sed -i "s/TOKENNAME/$NAME/g" ./contract.go