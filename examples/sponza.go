package main

// http://graphics.cs.williams.edu/data/meshes/dabrovic-sponza.zip

import (
	"flag"
	. "github.com/fogleman/pt/pt"
)

var pathTemplate = flag.String("path", "sponza-%04d.npy", "")

func main() {
	scene := Scene{}
	material := GlossyMaterial(HexColor(0xFCFAE1), 1.5, Radians(20))
	mesh, err := LoadOBJ("examples/dabrovic-sponza/sponza.obj", material)
	if err != nil {
		panic(err)
	}
	mesh.MoveTo(Vector{}, Vector{0.5, 0, 0.5})
	scene.Add(mesh)
	scene.Add(NewSphere(Vector{0, 20, 0}, 3, LightMaterial(White, 100)))
	camera := LookAt(Vector{-10, 2, 0}, Vector{0, 4, 0}, Vector{0, 1, 0}, 45)
	sampler := NewSampler(4, 4)
	// renderer := NewRenderer(&scene, &camera, sampler, 960, 540)
	renderer := NewRenderer(&scene, &camera, sampler, 400, 300)
	// renderer.IterativeRender("sponza-out%03d.png", 50)
	renderer.ExportFeatures(*pathTemplate, 1024)
}
