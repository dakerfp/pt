package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "teapot-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	scene.Add(NewSphere(Vector{-2, 5, -3}, 0.5, LightMaterial(White, 50)))
	scene.Add(NewSphere(Vector{5, 5, -3}, 0.5, LightMaterial(White, 50)))
	scene.Add(NewCube(Vector{-30, -1, -30}, Vector{30, 0, 30}, SpecularMaterial(HexColor(0xFCFAE1), 2)))
	mesh, err := LoadOBJ("examples/teapot.obj", SpecularMaterial(HexColor(0xB9121B), 2))
	if err != nil {
		panic(err)
	}
	scene.Add(mesh)
	camera := LookAt(Vector{2, 5, -6}, Vector{0.5, 1, 0}, Vector{0, 1, 0}, 45)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
