package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

// http://graphics.cs.williams.edu/data/meshes/dragon.zip

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "dragon-%04d.npy", "")

func main() {
	flag.Parse()
	scene := Scene{}

	material := GlossyMaterial(HexColor(0xB7CA79), 1.5, Radians(20))
	mesh, err := LoadOBJ("examples/dragon.obj", material)
	if err != nil {
		panic(err)
	}
	mesh.FitInside(Box{Vector{-1, 0, -1}, Vector{1, 2, 1}}, Vector{0.5, 0, 0.5})
	scene.Add(mesh)

	floor := GlossyMaterial(HexColor(0xD8CAA8), 1.2, Radians(5))
	scene.Add(NewCube(Vector{-50, -50, -50}, Vector{50, 0, 50}, floor))

	light := LightMaterial(White, 75)
	scene.Add(NewSphere(Vector{-1, 10, 0}, 1, light))

	mouth := LightMaterial(HexColor(0xFFFAD5), 500)
	scene.Add(NewSphere(V(-0.05, 1, -0.5), 0.03, mouth))

	camera := LookAt(Vector{-3, 2, -1}, Vector{0, 0.6, -0.1}, Vector{0, 1, 0}, 35)
	camera.SetFocus(Vector{0, 1, -0.5}, 0.03)
	sampler := NewSampler(4, 8)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
