#!/bin/bash

set -e

INV=$(pwd)/../involucro

cd ../examples/hello_world
$INV greet

exit 0
