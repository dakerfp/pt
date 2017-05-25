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
	material := GlossyMaterial(HexColor(0xF2F2F2), 1.5, Radians(20))
	scene.Add(NewCube(Vector{-100, -1, -100}, Vector{100, 0, 100}, material))
	heart := GlossyMaterial(HexColor(0xF60A20), 1.5, Radians(20))
	mesh, err := LoadSTL("examples/love.stl", heart)
	if err != nil {
		panic(err)
	}
	mesh.FitInside(Box{Vector{-0.5, 0, -0.5}, Vector{0.5, 1, 0.5}}, Vector{0.5, 0, 0.5})
	scene.Add(mesh)
	scene.Add(NewSphere(Vector{-2, 10, 2}, 1, LightMaterial(White, 30)))
	scene.Add(NewSphere(Vector{0, 10, 2}, 1, LightMaterial(White, 30)))
	scene.Add(NewSphere(Vector{2, 10, 2}, 1, LightMaterial(White, 30)))
	camera := LookAt(Vector{0, 1.5, 2}, Vector{0, 0.5, 0}, Vector{0, 1, 0}, 35)
	sampler := NewSampler(4, 4)
	renderer := NewRenderer(&scene, &camera, sampler, *width, *height, *spp)
	
	//renderer.IterativeRender("out%03d.png", 1000)
	renderer.ExportFeatures(*pathTemplate, *interactions)
}
