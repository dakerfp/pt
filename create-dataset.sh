#! /bin/bash

mkdir -p dataset

go build ./examples/cube.go
./cube -w 50 -h 60 -spp 1024 -path dataset/cube-1-%04d.npy
./cube -w 100 -h 100 -spp 1024 -path dataset/cube-2-%04d.npy
./cube -w 400 -h 300 -spp 1024 -path dataset/cube-3-%04d.npy
./cube -w 1200 -h 800 -spp 1024 -path dataset/cube-4-%04d.npy