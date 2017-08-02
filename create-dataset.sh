#! /bin/bash

SOURCE="${BASH_SOURCE[0]}"
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
# EXAMPLES="$(ls examples/ | grep .go | grep -v volume | grep -v dragon.obj | sed -e "s/.go$//")"
EXAMPLES="
cornell
go
dragon
cubes
teapot
hdri
maze
gopher
toybrick
counties
materials
ellipsoid
runway
example
love
bunny
craft
cylinder
volume
sdf"

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

function augment_multires()
{
	DIRNAME="./dataset/$1"
	ZIPFILE="./dataset/$1.zip"
	go build -o "./$1" "./examples/$1.go"
	"./$1" -w 32 -h 32 -interactions $2 -path "/tmp/1-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/1-%04d.npy" "1-%04d.npy" $2
	"./$1" -w 60 -h 40 -interactions $2 -path "/tmp/2-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/2-%04d.npy" "2-%04d.npy" $2
	"./$1" -w 120 -h 80 -interactions $2 -path "/tmp/3-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/3-%04d.npy" "3-%04d.npy" $2
	"./$1" -w 240 -h 160 -interactions $2 -path "/tmp/4-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/4-%04d.npy" "4-%04d.npy" $2
	"./$1" -w 640 -h 480 -interactions $2 -path "/tmp/5-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/5-%04d.npy" "5-%04d.npy" $2
	"./$1" -w 1280 -h 960 -interactions $2 -path "/tmp/6-%04d.npy"
	python ./update-zip.py "./dataset/$1.zip" "/tmp/6-%04d.npy" "6-%04d.npy" $2
	rm -rf /tmp/*.npy
}


# for EXAMPLE in $EXAMPLES
# do
# 	render_multires $EXAMPLE 1024
# done

for EXAMPLE in $EXAMPLES
do
	augment_multires $EXAMPLE 32
done



# "./$1" -w 60 -h 40 -interactions $2 -path "dataset/$1-1-%04d.npy"
# "./$1" -w 120 -h 80 -interactions $2 -path "dataset/$1-2-%04d.npy"
# "./$1" -w 240 -h 160 -interactions $2 -path "dataset/$1-3-%04d.npy"
# "./$1" -w 640 -h 480 -interactions $2 -path dataset/$1-4-%04d.npy
# "./$1" -w 1280 -h 960 -interactions $2 -path dataset/$1-5-%04d.npy
