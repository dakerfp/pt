#! /bin/bash

SOURCE="${BASH_SOURCE[0]}"
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
# EXAMPLES="$(ls examples/ | grep .go | grep -v volume | grep -v dragon.obj | sed -e "s/.go$//")"
EXAMPLES="veach_scene
volume"

cd "${DIR}"
mkdir -p dataset

function render_multires()
{
	DIRNAME="./dataset/$1"
	go build -o "./$1" "./examples/$1.go"
	mkdir -p "${DIRNAME}"
	"./$1" -w 32 -h 32 -interactions $2 -path "${DIRNAME}/1-%04d.npy"
	"./$1" -w 60 -h 40 -interactions $2 -path "${DIRNAME}/2-%04d.npy"
	"./$1" -w 120 -h 80 -interactions $2 -path "${DIRNAME}/3-%04d.npy"
	"./$1" -w 240 -h 160 -interactions $2 -path "${DIRNAME}/4-%04d.npy"
	"./$1" -w 640 -h 480 -interactions $2 -path "${DIRNAME}/5-%04d.npy"
	"./$1" -w 1280 -h 960 -interactions $2 -path "${DIRNAME}/6-%04d.npy"
	cd "./dataset/$1"
	ls | xargs zip "archive.zip"
	cd -
	mv "./dataset/$1/archive.zip" "./dataset/$1.zip"
	rm -rf "${DIRNAME}"

	rm "./$1"
}

for EXAMPLE in $EXAMPLES
do
	render_multires $EXAMPLE 1024
done

# "./$1" -w 60 -h 40 -interactions $2 -path "dataset/$1-1-%04d.npy"
# "./$1" -w 120 -h 80 -interactions $2 -path "dataset/$1-2-%04d.npy"
# "./$1" -w 240 -h 160 -interactions $2 -path "dataset/$1-3-%04d.npy"
# "./$1" -w 640 -h 480 -interactions $2 -path dataset/$1-4-%04d.npy
# "./$1" -w 1280 -h 960 -interactions $2 -path dataset/$1-5-%04d.npy
