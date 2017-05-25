package main

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var width = flag.Int("w", 500, "width")
var height = flag.Int("h", 300, "height")
var spp = flag.Int("spp", 1, "spp")
var interactions = flag.Int("interactions", 1024, "interactions")
var pathTemplate = flag.String("path", "hits-%04d.npy", "")

func main() {
	flag.Parse()

	scene := Scene{}
	material := DiffuseMaterial(Color{0.95, 0.95, 1})
	light := LightMaterial(White, 300)
	scene.Add(NewSphere(Vector{-0.75, -0.75, 5}, 0.25, light))
	scene.Add(NewCube(Vector{-1000, -1000, -1000}, Vector{1000, 1000, 0}, material))
	mesh, err := LoadSTL("examples/hits.stl", material)
	mesh.SmoothNormalsThreshold(Radians(10))
	mesh.FitInside(Box{V(-1, -1, 0), V(1, 1, 2)}, V(0.5, 0.5, 0))
	if err != nil {
		panic(err)
	}
	scene.Add(mesh)
	camera := LookAt(Vector{1.6, -3, 2}, Vector{-0.25, 0.5, 0.5}, Vector{0, 0, 1}, 50)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
